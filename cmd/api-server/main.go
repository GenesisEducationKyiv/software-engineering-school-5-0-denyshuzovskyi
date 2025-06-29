package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/emailclient"
	weatherprovider "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weather"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weather/weatherapi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weather/weatherstack"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/cron"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/database"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/logger"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/repository/postgresql"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/server/handler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/notification"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/weather"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/weatherupd"
	nimbusvalidator "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/validator"
	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mailgun/mailgun-go/v4"
)

func main() {
	cfg := config.ReadConfig("./config/config.yaml")
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	weatherLog := slog.New(slog.NewJSONHandler(logger.SetUpRotator("logs/weather.log"), &slog.HandlerOptions{Level: slog.LevelInfo}))

	if err := runApp(cfg, weatherLog, log); err != nil {
		log.Error("fatal error", "error", err)
		os.Exit(1)
	}
}

func runApp(cfg *config.Config, weatherLog *slog.Logger, log *slog.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	client := &http.Client{}
	weatherapiClient := weatherapi.NewClient(cfg.WeatherProvider.Url, cfg.WeatherProvider.Key, client, log)
	weatherapiProvider := weatherprovider.NewLoggingWeatherProvider("weatherapi.com", weatherapiClient, weatherLog, log)
	weatherstackClient := weatherstack.NewClient(cfg.FallbackWeatherProvider.Url, cfg.FallbackWeatherProvider.Key, client, log)
	weatherstackProvider := weatherprovider.NewLoggingWeatherProvider("weatherstack.com", weatherstackClient, weatherLog, log)
	chainWeatherProvider := weatherprovider.NewChainWeatherProvider(log, weatherapiProvider, weatherstackProvider)

	emailClient := emailclient.NewEmailClient(mailgun.NewMailgun(cfg.EmailService.Domain, cfg.EmailService.Key))

	weatherRepository := postgresql.NewWeatherRepository()
	subscriberRepository := postgresql.NewSubscriberRepository()
	subscriptionRepository := postgresql.NewSubscriptionRepository()
	tokenRepository := postgresql.NewTokenRepository()

	weatherService := weather.NewWeatherService(db, chainWeatherProvider, weatherRepository, log)
	subscriptionService := subscription.NewSubscriptionService(db, chainWeatherProvider, subscriberRepository, subscriptionRepository, tokenRepository, emailClient, cfg.Emails.ConfirmationEmail, cfg.Emails.ConfirmationSuccessfulEmail, cfg.Emails.UnsubscribeEmail, log)
	notificationService := notification.NewNotificationService(emailClient)
	weatherUpdateSendingService := weatherupd.NewWeatherUpdateSendingService(subscriptionService, weatherService, notificationService, log)

	validate := validator.New()
	locationValidator := nimbusvalidator.NewLocationValidator(validate)
	subscriptionValidator := nimbusvalidator.NewSubscriptionValidator(validate)

	weatherHandler := handler.NewWeatherHandler(weatherService, locationValidator, log)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, subscriptionValidator, log)

	cron, err := cron.SetUpCronJobs(ctx, weatherUpdateSendingService, cfg.Emails.WeatherEmail, log)
	if err != nil {
		return err
	}
	cron.Start()
	defer cron.Stop()

	mux := server.InitMux(weatherHandler, subscriptionHandler)
	srv := &http.Server{
		Addr:    net.JoinHostPort(cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler: mux,
	}

	log.Info("starting http server", "addr", srv.Addr)

	return srv.ListenAndServe()
}
