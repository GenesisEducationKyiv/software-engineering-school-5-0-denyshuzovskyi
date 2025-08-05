package service

import (
	"context"
	"fmt"
	"log/slog"

	v1 "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-proto/gen/go/notification/v1"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/dto"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/error"
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

func (s *EmailSendingService) SendConfirmation(ctx context.Context, req *v1.SendConfirmationRequest) (*v1.SendConfirmationResponse, error) {
	template := s.emailTemplates.Confirmation

	email := dto.SimpleEmail{
		From:    template.From,
		To:      req.GetNotificationWithToken().GetNotification().GetTo(),
		Subject: template.Subject,
		Text: fmt.Sprintf(
			template.Text,
			req.GetNotificationWithToken().GetToken(),
		),
	}

	if err := s.sender.Send(ctx, email); err != nil {
		s.log.Error("failed to send confirmation email", "error", err)
		return &v1.SendConfirmationResponse{}, commonerrors.ErrEmailSendingFailed
	}
	s.log.Debug("sent confirmation email successfully", "to", req.GetNotificationWithToken().GetNotification().GetTo())

	return &v1.SendConfirmationResponse{}, nil
}

func (s *EmailSendingService) SendConfirmationSuccess(ctx context.Context, req *v1.SendConfirmationSuccessRequest) (*v1.SendConfirmationSuccessResponse, error) {
	template := s.emailTemplates.ConfirmationSuccess

	email := dto.SimpleEmail{
		From:    template.From,
		To:      req.GetNotificationWithToken().GetNotification().GetTo(),
		Subject: template.Subject,
		Text: fmt.Sprintf(
			template.Text,
			req.GetNotificationWithToken().GetToken(),
		),
	}

	if err := s.sender.Send(ctx, email); err != nil {
		s.log.Error("failed to send confirmation success email", "error", err)
		return &v1.SendConfirmationSuccessResponse{}, commonerrors.ErrEmailSendingFailed
	}
	s.log.Debug("sent confirmation success email successfully", "to", req.GetNotificationWithToken().GetNotification().GetTo())

	return &v1.SendConfirmationSuccessResponse{}, nil
}

func (s *EmailSendingService) SendUnsubscribeSuccess(ctx context.Context, req *v1.SendUnsubscribeSuccessRequest) (*v1.SendUnsubscribeSuccessResponse, error) {
	template := s.emailTemplates.UnsubscribeSuccess

	email := dto.SimpleEmail{
		From:    template.From,
		To:      req.GetNotification().GetTo(),
		Subject: template.Subject,
		Text:    template.Text,
	}

	if err := s.sender.Send(ctx, email); err != nil {
		s.log.Error("failed to send unsubscribe success email", "error", err)
		return &v1.SendUnsubscribeSuccessResponse{}, commonerrors.ErrEmailSendingFailed
	}
	s.log.Debug("sent unsubscribe success email successfully", "to", req.GetNotification().GetTo())

	return &v1.SendUnsubscribeSuccessResponse{}, nil
}

func (s *EmailSendingService) SendWeatherUpdate(ctx context.Context, req *v1.SendWeatherUpdateRequest) (*v1.SendWeatherUpdateResponse, error) {
	template := s.emailTemplates.WeatherUpdate

	email := dto.SimpleEmail{
		From:    template.From,
		To:      req.GetWeatherUpdateNotification().GetNotificationWithToken().GetNotification().GetTo(),
		Subject: template.Subject,
		Text: fmt.Sprintf(
			template.Text,
			req.GetWeatherUpdateNotification().GetWeather().GetLocation(),
			req.GetWeatherUpdateNotification().GetWeather().GetTemperature(),
			req.GetWeatherUpdateNotification().GetWeather().GetHumidity(),
			req.GetWeatherUpdateNotification().GetWeather().GetDescription(),
			req.GetWeatherUpdateNotification().GetNotificationWithToken().GetToken(),
		),
	}

	if err := s.sender.Send(ctx, email); err != nil {
		s.log.Error("failed to send weather update email", "error", err)
		return &v1.SendWeatherUpdateResponse{}, commonerrors.ErrEmailSendingFailed
	}
	s.log.Debug("sent weather update email successfully", "to", req.GetWeatherUpdateNotification().GetNotificationWithToken().GetNotification().GetTo())

	return &v1.SendWeatherUpdateResponse{}, nil
}
