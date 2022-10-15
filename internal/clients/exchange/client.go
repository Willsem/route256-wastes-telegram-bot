package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/clients/exchange/dto"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/models"
)

type Config struct {
	Endpoint string `yaml:"endpoint"`
}

type Client struct {
	endpoint   *url.URL
	httpClient http.Client
}

func NewClient(config Config) (*Client, error) {
	return NewClientWithHttpClient(config, http.Client{})
}

func NewClientWithHttpClient(config Config, httpClient http.Client) (*Client, error) {
	endpoint, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint: %w", err)
	}

	return &Client{
		endpoint:   endpoint,
		httpClient: httpClient,
	}, nil
}

func (c *Client) GetExchange(ctx context.Context, base string, symbols []string) (*models.ExchangeData, error) {
	reqURL := c.endpoint
	values := reqURL.Query()
	values.Add("base", base)
	values.Add("symbols", strings.Join(symbols, ","))
	reqURL.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	var result dto.ExchangeData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("parse body: %w", err)
	}

	return models.NewExchangeData(result.Base, result.Rates), nil
}
