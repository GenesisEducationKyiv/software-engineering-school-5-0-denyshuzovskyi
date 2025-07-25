package main

import (
	"log/slog"
	"net"
	"os"

	v1 "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-proto/gen/go/notification/v1"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/client/emailclient"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/server"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/notification-srv/internal/service"
	"github.com/mailgun/mailgun-go/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.ReadConfig("./config/config.yaml")
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	emailClient := emailclient.NewEmailClient(mailgun.NewMailgun(cfg.EmailService.Domain, cfg.EmailService.Key))
	emailSendingService := service.NewEmailSendingService(cfg.EmailTemplates, emailClient, log)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	v1.RegisterNotificationServiceServer(grpcServer, server.NewNotificationServer(emailSendingService, log))

	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.GRPCServer.Host, cfg.GRPCServer.Port))
	if err != nil {
		log.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	log.Info("starting gRPC server on on", "addr", listener.Addr().String())
	if err := grpcServer.Serve(listener); err != nil {
		log.Error("failed to serve", "error", err)
	}
}
