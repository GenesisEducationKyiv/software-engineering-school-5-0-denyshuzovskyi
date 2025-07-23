package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/client/emailclient"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/consumer"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/service"
	"github.com/mailgun/mailgun-go/v4"
	amqp "github.com/rabbitmq/amqp091-go"
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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	conn, err := amqp.Dial(cfg.RabbitMQ.Url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", "error", err)
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Error("failed to close channel", "error", err)
		}
	}(ch)

	err = rabbitmq.SetUpQueue(ch, cfg.RabbitMQ.Exchange, cfg.RabbitMQ.Queue)
	if err != nil {
		return fmt.Errorf("failed to set up RabbitMQ queue: %w", err)
	}

	emailClient := emailclient.NewEmailClient(mailgun.NewMailgun(cfg.EmailService.Domain, cfg.EmailService.Key))
	emailSendingService := service.NewEmailSendingService(cfg.EmailTemplates, emailClient, log)
	notificationCommandConsumer := consumer.NewNotificationCommandConsumer(ch, cfg.RabbitMQ.Queue, emailSendingService, log)

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("started consuming notification commands")
		if err := notificationCommandConsumer.StartConsuming(ctx); err != nil && ctx.Err() == nil {
			errCh <- err
		}
	}()

	select {
	case sig := <-sigCh:
		log.Info("shutdown signal received", "signal", sig)
		cancel()

	case err := <-errCh:
		log.Error("consumer exited unexpectedly", "error", err)
		cancel()
	}

	wg.Wait()
	log.Info("graceful shutdown complete")

	return nil
}
