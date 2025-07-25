package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer              HTTPServer      `yaml:"server"`
	Datasource              Datasource      `yaml:"datasource"`
	WeatherProvider         WeatherProvider `yaml:"weather-provider" env-prefix:"WEATHER_PROVIDER_"`
	FallbackWeatherProvider WeatherProvider `yaml:"fallback-weather-provider" env-prefix:"FALLBACK_WEATHER_PROVIDER_"`
	Redis                   Redis           `yaml:"redis"`
	EmailService            EmailService    `yaml:"email-service"`
	Emails                  Emails          `yaml:"emails"`
}

type HTTPServer struct {
	Host string `yaml:"host" env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
}

type Datasource struct {
	Url string `yaml:"url" env:"DATABASE_URL"`
}

type WeatherProvider struct {
	Url string `yaml:"url" env:"URL"`
	Key string `yaml:"key" env:"KEY"`
}

type Redis struct {
	Url      string        `yaml:"url" env:"REDIS_URL"`
	Password string        `yaml:"password" env:"REDIS_PASSWORD"`
	TTL      time.Duration `yaml:"ttl"`
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

type Emails struct {
	ConfirmationEmail           EmailData `yaml:"confirmation-email"`
	ConfirmationSuccessfulEmail EmailData `yaml:"confirmation-successful-email"`
	WeatherEmail                EmailData `yaml:"weather-email"`
	UnsubscribeEmail            EmailData `yaml:"unsubscribe-email"`
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

	cfg.Emails.ConfirmationEmail.fillOutFromEmail(cfg.EmailService.Sender)
	cfg.Emails.ConfirmationSuccessfulEmail.fillOutFromEmail(cfg.EmailService.Sender)
	cfg.Emails.WeatherEmail.fillOutFromEmail(cfg.EmailService.Sender)
	cfg.Emails.UnsubscribeEmail.fillOutFromEmail(cfg.EmailService.Sender)

	return &cfg
}
