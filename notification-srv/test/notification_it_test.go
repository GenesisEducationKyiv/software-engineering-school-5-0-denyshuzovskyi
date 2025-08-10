//go:build integration

package test

import (
	"context"
	"errors"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	ptestutil "github.com/prometheus/client_golang/prometheus/testutil"
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/command/notification"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/consumer"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	tcrabbitmq "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

var notificationCommands = []notification.NotificationCommand{
	&notification.SendConfirmation{
		NotificationWithToken: notification.NotificationWithToken{
			Notification: notification.Notification{
				To: "test@test.com",
			},
			Token: "test_token",
		},
	},
	&notification.SendConfirmationSuccess{
		NotificationWithToken: notification.NotificationWithToken{
			Notification: notification.Notification{
				To: "test@test.com",
			},
			Token: "test_token",
		},
	},
	&notification.SendWeatherUpdate{
		NotificationWithToken: notification.NotificationWithToken{
			Notification: notification.Notification{
				To: "test@test.com",
			},
			Token: "test_token",
		},
		Weather: notification.Weather{
			Location:    "Kyiv",
			Temperature: float32(22),
			Humidity:    float32(40),
			Description: "Clear",
		},
	},
	&notification.SendUnsubscribeSuccess{
		Notification: notification.Notification{
			To: "test@test.com",
		},
	},
}

func SetUpRabbitMQ(t *testing.T) *amqp.Channel {
	t.Helper()
	ctx := context.Background()

	rabbitmqContainer, err := tcrabbitmq.Run(ctx, "rabbitmq:4-alpine")
	require.NoError(t, err)

	amqpURL, err := rabbitmqContainer.AmqpURL(ctx)
	require.NoError(t, err)
	t.Logf("RabbitMQ amqpURL: %s", amqpURL)

	rabbitmqRes, err := rabbitmq.InitRabbitMQ(amqpURL)
	require.NoError(t, err)

	t.Cleanup(func() {
		t.Log("started cleanup")

		// closing channel before connection might freeze the cleanup
		require.NoError(t, rabbitmqRes.Close())
		require.NoError(t, rabbitmqContainer.Terminate(ctx))

		t.Log("finished cleanup")
	})

	return rabbitmqRes.Channel
}

func TestNotificationCommandConsumer_HandleAllNotificationTypes(t *testing.T) {
	ch := SetUpRabbitMQ(t)

	defaultExchange := ""
	queueName := "test-queue"
	_, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	require.NoError(t, err)

	// publishing to exchange = "" (default exchange) with routingKey = queueName routes directly to a queue with that name, if it exists
	routingKey := rabbitmq.RoutingKey(queueName)

	publisher := rabbitmq.NewPublisher(ch, defaultExchange)

	cfg := config.ReadConfig("./../config/config.yaml")
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	emailSenderMock := service.NewMockEmailSender(t)
	emailSent := make(chan struct{}, len(notificationCommands))

	emailSentMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_email_sent_total",
		Help: "",
	}, []string{"type"})

	emailFailedMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_email_failed_total",
		Help: "",
	}, []string{"type"})

	emailMetrics := metrics.NewPrometheusEmailMetrics(
		metrics.WithEmailSentCounterVec(emailSentMetric),
		metrics.WithEmailFailedCounterVec(emailFailedMetric),
	)
	emailTypes := []notification.CommandType{notification.Confirmation, notification.ConfirmationSuccess, notification.WeatherUpdate, notification.UnsubscribeSuccess}
	emailTypesStr := make([]string, len(emailTypes))
	for i := range len(emailTypes) {
		emailTypesStr[i] = string(emailTypes[i])
	}
	emailMetrics.Init(emailTypesStr)

	emailSendingService := service.NewEmailSendingService(cfg.EmailTemplates, emailSenderMock, emailMetrics, log)
	notificationCommandDispatcher := consumer.NewNotificationCommandDispatcher(emailSendingService)
	notificationCommandConsumer := consumer.NewNotificationCommandConsumer(ch, queueName, notificationCommandDispatcher, log)

	var sentEmails []dto.SimpleEmail
	emailSenderMock.EXPECT().
		Send(mock.Anything, mock.AnythingOfType("dto.SimpleEmail")).
		Run(func(ctx context.Context, email dto.SimpleEmail) {
			sentEmails = append(sentEmails, email)
			emailSent <- struct{}{}
		}).
		Return(nil).Times(len(notificationCommands))

	ctx, cancel := context.WithCancel(t.Context())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("starting consumer")
		err := notificationCommandConsumer.StartConsuming(ctx)
		if !errors.Is(err, context.Canceled) {
			t.Logf("consumer exited: %v", err)
			require.NoError(t, err)
		}
		t.Log("consumer exited")
	}()

	for _, command := range notificationCommands {
		publishNotificationCommand(t, publisher, routingKey, command)
	}

	for i := range len(notificationCommands) {
		select {
		case _ = <-emailSent:
			t.Logf("email sent successfully %d", i)
		}
	}

	t.Log("cancelling context to stop consumer")
	cancel()

	require.Len(t, sentEmails, 4)
	require.Equal(t, cfg.EmailTemplates.Confirmation.Subject, sentEmails[0].Subject)
	require.Equal(t, cfg.EmailTemplates.ConfirmationSuccess.Subject, sentEmails[1].Subject)
	require.Equal(t, cfg.EmailTemplates.WeatherUpdate.Subject, sentEmails[2].Subject)
	require.Equal(t, cfg.EmailTemplates.UnsubscribeSuccess.Subject, sentEmails[3].Subject)

	delta := 0.1
	for _, emailTypeStr := range emailTypesStr {
		assertEmailMetrics(t, emailSentMetric, emailFailedMetric, emailTypeStr, 1, 0, delta)
	}
	assertEmailMetrics(t, emailSentMetric, emailFailedMetric, "unknown", 0, 0, delta)

	wg.Wait()
	t.Log("test completed successfully")
}

func publishNotificationCommand(t *testing.T, publisher *rabbitmq.Publisher, routingKey rabbitmq.RoutingKey, cmd notification.NotificationCommand) {
	t.Helper()

	bytes, err := notification.MarshalEnvelopeFromCommand(cmd)
	require.NoError(t, err)

	err = publisher.Publish(t.Context(), routingKey, bytes)
	require.NoError(t, err)

	t.Logf("published notification %s", cmd.Type())
}

func assertEmailMetrics(
	t *testing.T,
	sent *prometheus.CounterVec,
	failed *prometheus.CounterVec,
	emailType string,
	expectedSent int,
	expectedFailed int,
	delta float64,
) {
	t.Helper()

	require.InDelta(t, float64(expectedSent), ptestutil.ToFloat64(sent.WithLabelValues(emailType)), delta, "emails sent mismatch")
	require.InDelta(t, float64(expectedFailed), ptestutil.ToFloat64(failed.WithLabelValues(emailType)), delta, "emails failed mismatch")
}
