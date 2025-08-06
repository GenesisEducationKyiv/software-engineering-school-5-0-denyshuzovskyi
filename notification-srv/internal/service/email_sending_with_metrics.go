package service

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/command/notification"
)

type EmailMetrics interface {
	RecordEmailSent(string)
	RecordEmailFailed(string)
}

type EmailNotificationSendingService interface {
	SendConfirmation(context.Context, notification.SendConfirmation) error
	SendConfirmationSuccess(context.Context, notification.SendConfirmationSuccess) error
	SendWeatherUpdate(context.Context, notification.SendWeatherUpdate) error
	SendUnsubscribeSuccess(context.Context, notification.SendUnsubscribeSuccess) error
}

type EmailSendingServiceWithMetrics struct {
	emailSendingService EmailNotificationSendingService
	emailMetrics        EmailMetrics
}

func NewEmailSendingServiceWithMetrics(service EmailNotificationSendingService, emailMetrics EmailMetrics) *EmailSendingServiceWithMetrics {
	return &EmailSendingServiceWithMetrics{
		emailSendingService: service,
		emailMetrics:        emailMetrics,
	}
}

func (s *EmailSendingServiceWithMetrics) SendConfirmation(ctx context.Context, confirmation notification.SendConfirmation) error {
	err := s.emailSendingService.SendConfirmation(ctx, confirmation)

	return s.recordMetric(err, &confirmation)
}

func (s *EmailSendingServiceWithMetrics) SendConfirmationSuccess(ctx context.Context, confirmationSuccess notification.SendConfirmationSuccess) error {
	err := s.emailSendingService.SendConfirmationSuccess(ctx, confirmationSuccess)

	return s.recordMetric(err, &confirmationSuccess)
}

func (s *EmailSendingServiceWithMetrics) SendWeatherUpdate(ctx context.Context, weatherUpdate notification.SendWeatherUpdate) error {
	err := s.emailSendingService.SendWeatherUpdate(ctx, weatherUpdate)

	return s.recordMetric(err, &weatherUpdate)
}

func (s *EmailSendingServiceWithMetrics) SendUnsubscribeSuccess(ctx context.Context, unsubscribeSuccess notification.SendUnsubscribeSuccess) error {
	err := s.emailSendingService.SendUnsubscribeSuccess(ctx, unsubscribeSuccess)

	return s.recordMetric(err, &unsubscribeSuccess)
}

func (s *EmailSendingServiceWithMetrics) recordMetric(err error, command notification.NotificationCommand) error {
	if err != nil {
		s.emailMetrics.RecordEmailFailed(command.Type())
		return err
	}

	s.emailMetrics.RecordEmailSent(command.Type())
	return nil
}
