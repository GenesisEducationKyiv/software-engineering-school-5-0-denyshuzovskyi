package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/emailclient"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weatherapi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config/email"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/cron"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/database"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/repository/postgresql"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/server/handler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/notification"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/weather"
	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mailgun/mailgun-go/v4"
)

func main() {
	cfg := config.ReadConfig("./config/config.yaml")
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if err := runApp(cfg, log); err != nil {
		log.Error("fatal error", "error", err)
		os.Exit(1)
	}
}

func runApp(cfg *config.Config, log *slog.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	validate := validator.New()

	db, err := database.InitDB(ctx, cfg.Datasource.Url, log)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Error("unable to close connection pool", "error", err)
		}
	}(db)

	if err := database.RunMigrations(db, ".", log); err != nil {
		return err
	}

	emailDataMap := email.PrepareEmailData(cfg)
	confirmEmailData, confOk := emailDataMap["confirmation"]
	confirmSuccessEmailData, confSuccessOk := emailDataMap["confirmation-successful"]
	weatherEmailData, weatherOk := emailDataMap["weather"]
	unsubEmailData, unsubOk := emailDataMap["unsubscribe"]
	if !confOk || !confSuccessOk || !weatherOk || !unsubOk {
		log.Error("cannot prepare email data")
		return commonerrors.ErrInvalidEmailData
	}

	weatherApiClient := weatherapi.NewClient(cfg.WeatherProvider.Url, cfg.WeatherProvider.Key, &http.Client{}, log)
	emailClient := emailclient.NewEmailClient(mailgun.NewMailgun(cfg.EmailService.Domain, cfg.EmailService.Key))

	weatherRepository := postgresql.NewWeatherRepository()
	subscriberRepository := postgresql.NewSubscriberRepository()
	subscriptionRepository := postgresql.NewSubscriptionRepository()
	tokenRepository := postgresql.NewTokenRepository()

	weatherService := weather.NewWeatherService(db, weatherApiClient, weatherRepository, log)
	subscriptionService := subscription.NewSubscriptionService(db, weatherApiClient, subscriberRepository, subscriptionRepository, tokenRepository, emailClient, confirmEmailData, confirmSuccessEmailData, unsubEmailData, log)
	notificationService := notification.NewNotificationService(db, weatherApiClient, weatherRepository, subscriberRepository, subscriptionRepository, tokenRepository, emailClient, log)

	weatherHandler := handler.NewWeatherHandler(weatherService, log)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, validate, log)

	if err = cron.StartCronJobs(ctx, notificationService, weatherEmailData, log); err != nil {
		return err
	}

	router := server.InitRouter(weatherHandler, subscriptionHandler)
	srv := &http.Server{
		Addr:    net.JoinHostPort(cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler: router,
	}

	log.Info("starting http server", "addr", srv.Addr)

	return srv.ListenAndServe()
}
