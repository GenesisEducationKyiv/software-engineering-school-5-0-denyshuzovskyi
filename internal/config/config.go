package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	HTTPServer      `yaml:"server"`
	Datasource      `yaml:"datasource"`
	WeatherProvider `yaml:"weather-provider"`
	EmailService    `yaml:"email-service"`
	Emails          []EmailData `yaml:"emails"`
}

type HTTPServer struct {
	Host string `yaml:"host" env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
}

type Datasource struct {
	Url string `yaml:"url" env:"DATABASE_URL"`
}

type WeatherProvider struct {
	Url string `yaml:"url" env:"WEATHER_PROVIDER_URL"`
	Key string `yaml:"key" env:"WEATHER_PROVIDER_KEY"`
}

type EmailService struct {
	Domain string `yaml:"domain" env:"EMAIL_SERVICE_DOMAIN"`
	Key    string `yaml:"key" env:"EMAIL_SERVICE_KEY"`
	Sender string `yaml:"sender"`
}

type EmailData struct {
	Name    string `yaml:"name"`
	Subject string `yaml:"subject"`
	Text    string `yaml:"text"`
	From    string
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

	return &cfg
}
