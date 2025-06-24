package weatherstack

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
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

func (c *Client) GetCurrentWeather(ctx context.Context, location string) (*dto.WeatherWithLocationDTO, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	q.Set("access_key", c.apiKey)
	q.Set("query", location)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.log.Error("failed to close body", "error", err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		return nil, commonerrors.ErrLocationNotFound
	default:
		return nil, fmt.Errorf("%w status code: %d", commonerrors.ErrUnexpectedStatusCode, resp.StatusCode)
	}

	var currentWeather CurrentWeather
	if err := json.NewDecoder(resp.Body).Decode(&currentWeather); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var weatherWithLocationDTO dto.WeatherWithLocationDTO
	weatherWithLocationDTO, err = CurrentWeatherToWeatherWithLocationDTO(currentWeather)
	if err != nil {
		return nil, fmt.Errorf("mapping error: %w", err)
	}
	return &weatherWithLocationDTO, nil
}
