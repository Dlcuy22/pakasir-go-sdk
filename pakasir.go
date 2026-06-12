package pakasir

import (
	"context"
	"errors"
	"net/http"
	"sync"
)

type Client struct {
	cfg        Config
	httpClient *http.Client
	watchers   map[string]context.CancelFunc
	mu         sync.Mutex
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.Slug == "" || cfg.APIKey == "" {
		return nil, errors.New("slug and API key are required")
	}

	return &Client{
		cfg:        cfg,
		httpClient: http.DefaultClient,
		watchers:   make(map[string]context.CancelFunc),
	}, nil
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}
