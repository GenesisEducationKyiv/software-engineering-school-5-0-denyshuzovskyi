package server

import (
	"context"
	"log/slog"

	v1 "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-proto/gen/go/notification/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationSender interface {
	SendConfirmation(ctx context.Context, req *v1.SendConfirmationRequest) (*v1.SendConfirmationResponse, error)
	SendConfirmationSuccess(ctx context.Context, req *v1.SendConfirmationSuccessRequest) (*v1.SendConfirmationSuccessResponse, error)
	SendUnsubscribeSuccess(ctx context.Context, req *v1.SendUnsubscribeSuccessRequest) (*v1.SendUnsubscribeSuccessResponse, error)
	SendWeatherUpdate(ctx context.Context, req *v1.SendWeatherUpdateRequest) (*v1.SendWeatherUpdateResponse, error)
}

type NotificationServer struct {
	v1.UnimplementedNotificationServiceServer
	sender NotificationSender
	log    *slog.Logger
}

func NewNotificationServer(sender NotificationSender, log *slog.Logger) *NotificationServer {
	return &NotificationServer{
		sender: sender,
		log:    log,
	}
}

func (s *NotificationServer) SendConfirmation(ctx context.Context, req *v1.SendConfirmationRequest) (*v1.SendConfirmationResponse, error) {
	sendConfirmResp, err := s.sender.SendConfirmation(ctx, req)
	if err != nil {
		return &v1.SendConfirmationResponse{}, status.Error(codes.Internal, err.Error())
	}

	return sendConfirmResp, nil
}

func (s *NotificationServer) SendConfirmationSuccess(ctx context.Context, req *v1.SendConfirmationSuccessRequest) (*v1.SendConfirmationSuccessResponse, error) {
	sendConfirmSuccessResp, err := s.sender.SendConfirmationSuccess(ctx, req)
	if err != nil {
		return &v1.SendConfirmationSuccessResponse{}, status.Error(codes.Internal, err.Error())
	}

	return sendConfirmSuccessResp, nil
}

func (s *NotificationServer) SendUnsubscribeSuccess(ctx context.Context, req *v1.SendUnsubscribeSuccessRequest) (*v1.SendUnsubscribeSuccessResponse, error) {
	sendUnsubSuccessResp, err := s.sender.SendUnsubscribeSuccess(ctx, req)
	if err != nil {
		return &v1.SendUnsubscribeSuccessResponse{}, status.Error(codes.Internal, err.Error())
	}

	return sendUnsubSuccessResp, nil
}

func (s *NotificationServer) SendWeatherUpdate(ctx context.Context, req *v1.SendWeatherUpdateRequest) (*v1.SendWeatherUpdateResponse, error) {
	sendWeatherUpdResp, err := s.sender.SendWeatherUpdate(ctx, req)
	if err != nil {
		return &v1.SendWeatherUpdateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return sendWeatherUpdResp, nil
}
