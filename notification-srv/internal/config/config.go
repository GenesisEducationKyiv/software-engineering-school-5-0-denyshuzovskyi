package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCServer     GRPCServer     `yaml:"grpc-server"`
	EmailService   EmailService   `yaml:"email-service"`
	EmailTemplates EmailTemplates `yaml:"email-templates"`
}

type GRPCServer struct {
	Host string `yaml:"host" env:"GRPC_SERVER_HOST"`
	Port string `yaml:"port" env:"GRPC_SERVER_PORT"`
}

type EmailService struct {
	Domain string `yaml:"domain" env:"EMAIL_SERVICE_DOMAIN"`
	Key    string `yaml:"key" env:"EMAIL_SERVICE_KEY"`
	Sender string `yaml:"sender"`
}

type EmailData struct {
	Subject string `yaml:"subject"`
	Text    string `yaml:"text"`
	From    string
}

func (ed *EmailData) fillOutFromEmail(from string) {
	ed.From = from
}

type EmailTemplates struct {
	Confirmation        EmailData `yaml:"confirmation"`
	ConfirmationSuccess EmailData `yaml:"confirmation-success"`
	WeatherUpdate       EmailData `yaml:"weather-update"`
	UnsubscribeSuccess  EmailData `yaml:"unsubscribe-success"`
}

func ReadConfig(configPath string) *Config {
	if configPath == "" {
		log.Fatal("configPath is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	cfg.EmailTemplates.Confirmation.fillOutFromEmail(cfg.EmailService.Sender)
	cfg.EmailTemplates.ConfirmationSuccess.fillOutFromEmail(cfg.EmailService.Sender)
	cfg.EmailTemplates.WeatherUpdate.fillOutFromEmail(cfg.EmailService.Sender)
	cfg.EmailTemplates.UnsubscribeSuccess.fillOutFromEmail(cfg.EmailService.Sender)

	return &cfg
}
