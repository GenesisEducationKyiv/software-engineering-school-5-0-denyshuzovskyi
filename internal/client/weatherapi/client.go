package weatherapi

import (
	"context"
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

func (c *Client) GetCurrentWeather(ctx context.Context, location string) (*model.Weather, error) {
	u, err := url.Parse(c.baseURL + "/current.json")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("q", location)
	q.Set("aqi", "no")
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

	weather := CurrentWeatherToWeather(currentWeather)
	return &weather, nil
}
