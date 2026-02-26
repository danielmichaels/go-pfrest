package pfrest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielmichaels/go-pfrest/api"
)

type Config struct {
	BaseURL            string
	InsecureSkipVerify bool
	HTTPClient         *http.Client

	BasicAuth *BasicAuthConfig
	APIKey    string
	JWTToken  string
}

type BasicAuthConfig struct {
	Username string
	Password string
}

type Client struct {
	raw *api.ClientWithResponses
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, errors.New("pfrest: BaseURL is required")
	}
	if !strings.HasPrefix(cfg.BaseURL, "http://") && !strings.HasPrefix(cfg.BaseURL, "https://") {
		cfg.BaseURL = "https://" + cfg.BaseURL
	}
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")

	authFn, err := authEditor(cfg)
	if err != nil {
		return nil, fmt.Errorf("pfrest: %w", err)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = tlsClient(cfg.InsecureSkipVerify)
	}

	opts := []api.ClientOption{
		api.WithHTTPClient(httpClient),
		api.WithRequestEditorFn(authFn),
	}

	raw, err := api.NewClientWithResponses(cfg.BaseURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("pfrest: create client: %w", err)
	}

	return &Client{raw: raw}, nil
}

func (c *Client) Raw() *api.ClientWithResponses {
	return c.raw
}
