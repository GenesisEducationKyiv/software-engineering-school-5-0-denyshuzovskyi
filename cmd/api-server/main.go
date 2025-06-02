package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/emailclient"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/client/weatherapi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/handler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/repository/posgresql"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/migrations"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/robfig/cron/v3"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.ReadConfig("./config/config.yaml")
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	validate := validator.New()

	db, err := sql.Open("pgx", cfg.Datasource.Url)
	if err != nil {
		log.Error("unable to open database", "error", err)
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Error("unable to close connection pool", "error", err)
		}
	}(db)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error("unable to acquire database driver", "error", err)
		os.Exit(1)
	}

	d, err := iofs.New(migrations.Files, ".")
	if err != nil {
		log.Error("unable to set up driver for io/fs#FS", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Error("unable to set up migrations", "error", err)
		os.Exit(1)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("all migrations have already been applied")
		} else {
			log.Error("unable to apply migrations", "error", err)
			os.Exit(1)
		}
	} else {
		log.Info("migration completed successfully")
	}

	emailDataMap := prepareEmailData(cfg)
	confirmEmailData, confOk := emailDataMap["confirmation"]
	confirmSuccessEmailData, confSuccessOk := emailDataMap["confirmation-successful"]
	weatherEmailData, weatherOk := emailDataMap["weather"]
	unsubEmailData, unsubOk := emailDataMap["unsubscribe"]
	if !confOk || !confSuccessOk || !weatherOk || !unsubOk {
		log.Error("cannot prepare email data")
		os.Exit(1)
	}

	weatherApiClient := weatherapi.NewClient(cfg.WeatherProvider.Url, cfg.WeatherProvider.Key, &http.Client{}, log)
	emailClient := emailclient.NewEmailClient(mailgun.NewMailgun(cfg.EmailService.Domain, cfg.EmailService.Key))
	locationRepository := posgresql.NewLocationRepository()
	weatherRepository := posgresql.NewWeatherRepository()
	subscriberRepository := posgresql.NewSubscriberRepository()
	subscriptionRepository := posgresql.NewSubscriptionRepository()
	tokenRepository := posgresql.NewTokenRepository()
	weatherService := service.NewWeatherService(db, weatherApiClient, locationRepository, weatherRepository, log)
	subscriptionService := service.NewSubscriptionService(db, weatherApiClient, locationRepository, subscriberRepository, subscriptionRepository, tokenRepository, emailClient, confirmEmailData, confirmSuccessEmailData, unsubEmailData, log)
	notificationService := service.NewNotificationService(db, weatherApiClient, locationRepository, weatherRepository, subscriberRepository, subscriptionRepository, tokenRepository, emailClient, log)
	weatherHandler := handler.NewWeatherHandler(weatherService, log)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, validate, log)

	c := cron.New()
	// daily 09:00
	_, err = c.AddFunc("0 9 * * *", func() {
		notificationService.SendDailyNotifications(weatherEmailData)
	})
	if err != nil {
		log.Error("failed to schedule notification service", "error", err)
		os.Exit(1)
	}
	// hourly
	_, err = c.AddFunc("0 * * * *", func() {
		notificationService.SendHourlyNotifications(weatherEmailData)
	})
	if err != nil {
		log.Error("failed to schedule notification service", "error", err)
		os.Exit(1)
	}
	c.Start()

	router := http.NewServeMux()
	router.HandleFunc("GET /weather", weatherHandler.GetCurrentWeather)
	router.HandleFunc("POST /subscribe", subscriptionHandler.Subscribe)
	router.HandleFunc("GET /confirm/{token}", subscriptionHandler.Confirm)
	router.HandleFunc("GET /unsubscribe/{token}", subscriptionHandler.Unsubscribe)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler: router,
	}

	log.Info("starting server", "host", cfg.HTTPServer.Host, "port", cfg.HTTPServer.Port)

	err = server.ListenAndServe()
	if err != nil {
		log.Error("failed to start server", "error", err)
		return
	}
}

func prepareEmailData(cfg *config.Config) map[string]config.EmailData {
	emailDataMap := make(map[string]config.EmailData)

	for _, email := range cfg.Emails {
		memail := email
		memail.From = cfg.EmailService.Sender
		emailDataMap[memail.Name] = memail
	}

	return emailDataMap
}
