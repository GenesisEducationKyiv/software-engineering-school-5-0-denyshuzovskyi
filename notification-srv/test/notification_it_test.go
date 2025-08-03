//go:build integration

package test

import (
	"context"
	"errors"
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
	"log/slog"
	"os"
	"sync"
	"testing"
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

	queueName := "test-queue"
	_, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	require.NoError(t, err)

	cfg := config.ReadConfig("./../config/config.yaml")
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	emailSenderMock := service.NewMockEmailSender(t)
	emailSent := make(chan struct{}, len(notificationCommands))

	emailSendingService := service.NewEmailSendingService(cfg.EmailTemplates, emailSenderMock, log)
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
		publishNotificationCommand(t, ch, queueName, command)
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

	wg.Wait()
	t.Log("test completed successfully")
}

func publishNotificationCommand(t *testing.T, ch *amqp.Channel, queue string, cmd notification.NotificationCommand) {
	t.Helper()

	bytes, err := notification.MarshalEnvelopeFromCommand(cmd)
	require.NoError(t, err)

	rabbitmqPublisher := rabbitmq.NewPublisher(ch, "")
	err = rabbitmqPublisher.Publish(t.Context(), queue, bytes)
	require.NoError(t, err)

	t.Logf("published notification %s", cmd.Type())
}
