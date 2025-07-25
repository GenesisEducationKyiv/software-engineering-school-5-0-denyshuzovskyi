package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/notification/command"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/dto"
)

type EmailSender interface {
	Send(context.Context, dto.SimpleEmail) error
}

type EmailSendingService struct {
	emailTemplates config.EmailTemplates
	sender         EmailSender
	log            *slog.Logger
}

func NewEmailSendingService(emailTemplates config.EmailTemplates, sender EmailSender, log *slog.Logger) *EmailSendingService {
	return &EmailSendingService{
		emailTemplates: emailTemplates,
		sender:         sender,
		log:            log,
	}
}

func (s *EmailSendingService) SendConfirmation(ctx context.Context, sendConfirm command.SendConfirmation) error {
	template := s.emailTemplates.Confirmation

	email := dto.SimpleEmail{
		From:    template.From,
		To:      sendConfirm.To,
		Subject: template.Subject,
		Text: fmt.Sprintf(
			template.Text,
			sendConfirm.Token,
		),
	}

	if err := s.sender.Send(ctx, email); err != nil {
		return fmt.Errorf("sending confirmation email: %w", err)
	}
	s.log.Debug("sent confirmation email successfully", "to", sendConfirm.To)

	return nil
}

func (s *EmailSendingService) SendConfirmationSuccess(ctx context.Context, sendConfirmSuccess command.SendConfirmationSuccess) error {
	template := s.emailTemplates.ConfirmationSuccess

	email := dto.SimpleEmail{
		From:    template.From,
		To:      sendConfirmSuccess.To,
		Subject: template.Subject,
		Text: fmt.Sprintf(
			template.Text,
			sendConfirmSuccess.Token,
		),
	}

	if err := s.sender.Send(ctx, email); err != nil {
		return fmt.Errorf("sending confirmation success email: %w", err)
	}
	s.log.Debug("sent confirmation success email successfully", "to", sendConfirmSuccess.To)

	return nil
}

func (s *EmailSendingService) SendWeatherUpdate(ctx context.Context, sendWeatherUpd command.SendWeatherUpdate) error {
	template := s.emailTemplates.WeatherUpdate

	email := dto.SimpleEmail{
		From:    template.From,
		To:      sendWeatherUpd.To,
		Subject: template.Subject,
		Text: fmt.Sprintf(
			template.Text,
			sendWeatherUpd.Weather.Location,
			sendWeatherUpd.Weather.Temperature,
			sendWeatherUpd.Weather.Humidity,
			sendWeatherUpd.Weather.Description,
			sendWeatherUpd.Token,
		),
	}

	if err := s.sender.Send(ctx, email); err != nil {
		return fmt.Errorf("sending weather update email: %w", err)
	}
	s.log.Debug("sent weather update email successfully", "to", sendWeatherUpd.To)

	return nil
}

func (s *EmailSendingService) SendUnsubscribeSuccess(ctx context.Context, unsubSuccess command.SendUnsubscribeSuccess) error {
	template := s.emailTemplates.UnsubscribeSuccess

	email := dto.SimpleEmail{
		From:    template.From,
		To:      unsubSuccess.To,
		Subject: template.Subject,
		Text:    template.Text,
	}

	if err := s.sender.Send(ctx, email); err != nil {
		return fmt.Errorf("sending unsubscribe email: %w", err)
	}
	s.log.Debug("sent unsubscribe success email successfully", "to", unsubSuccess.To)

	return nil
}
