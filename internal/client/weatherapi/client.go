package weatherapi

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
	log     *slog.Logger
}

func NewClient(baseURL, apiKey string, client *http.Client, log *slog.Logger) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  client,
		log:     log,
	}
}

func (c *Client) GetCurrentWeather(location string) (*model.WeatherWithLocation, error) {
	u, err := url.Parse(c.baseURL + "/current.json")
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %w", err)
	}

	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("q", location)
	q.Set("aqi", "no")
	u.RawQuery = q.Encode()

	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("failed to close body", "error", err)
		}
	}()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, commonerrors.ErrLocationNotFound
	}

	var weather CurrentWeather
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("failed decode response %w", err)
	}
	weatherWithLocation := CurrentWeatherToWeatherWithLocation(weather)

	return &weatherWithLocation, nil
}
