package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"net/http"
	"os"

	v1 "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-proto/gen/go/notification/v1"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/cache"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/client/notificationservice"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/client/weatherapi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/client/weatherstack"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/cron"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/database"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/logger"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/metrics"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/repository/postgresql"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/server/handler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/notification"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/weather"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/weatherprovider"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/service/weatherupd"
	validators "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/validator"
	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	metrics.Init()

	client := &http.Client{}
	weatherapiClient := weatherapi.NewClient(cfg.WeatherProvider.Url, cfg.WeatherProvider.Key, client, log)
	weatherapiProvider := weatherprovider.NewLoggingWeatherProvider("weatherapi.com", weatherprovider.NewWeatherapiProvider(weatherapiClient), weatherLog, log)
	weatherstackClient := weatherstack.NewClient(cfg.FallbackWeatherProvider.Url, cfg.FallbackWeatherProvider.Key, client, log)
	weatherstackProvider := weatherprovider.NewLoggingWeatherProvider("weatherstack.com", weatherprovider.NewWeatherstackProvider(weatherstackClient), weatherLog, log)
	chainWeatherProvider := weatherprovider.NewChainWeatherProvider(log, weatherapiProvider, weatherstackProvider)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Url,
		Password: cfg.Redis.Password,
	})
	redisCache := cache.NewRedisCache(redisClient)
	weatherCache := cache.NewJSONCache[dto.WeatherWithLocationDTO](redisCache, cfg.Redis.TTL)
	cachingWeatherProvider := weatherprovider.NewCachingWeatherProvider(weatherCache, chainWeatherProvider, metrics.WeatherCacheRequests, log)

	subscriberRepository := postgresql.NewSubscriberRepository()
	subscriptionRepository := postgresql.NewSubscriptionRepository()
	tokenRepository := postgresql.NewTokenRepository()

	conn, err := grpc.NewClient(cfg.NotificationService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Error("failed to close client", "error", err)
		}
	}(conn)
	notificationServiceClient := notificationservice.NewNotificationServiceClient(v1.NewNotificationServiceClient(conn))

	weatherService := weather.NewWeatherService(cachingWeatherProvider, log)
	subscriptionService := subscription.NewSubscriptionService(db, cachingWeatherProvider, subscriberRepository, subscriptionRepository, tokenRepository, notificationServiceClient, log)
	notificationService := notification.NewNotificationService(notificationServiceClient)
	weatherUpdateSendingService := weatherupd.NewWeatherUpdateSendingService(subscriptionService, weatherService, notificationService, log)

	validate := validator.New()
	locationValidator := validators.NewLocationValidator(validate)
	subscriptionValidator := validators.NewSubscriptionValidator(validate)

	weatherHandler := handler.NewWeatherHandler(weatherService, locationValidator, log)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, subscriptionValidator, log)

	cron, err := cron.SetUpCronJobs(ctx, weatherUpdateSendingService, log)
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
