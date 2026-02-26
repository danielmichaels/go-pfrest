package pfrest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielmichaels/go-pfrest/api"
)

func authEditor(cfg Config) (api.RequestEditorFn, error) {
	count := authMethodCount(cfg)
	if count == 0 {
		return nil, errors.New("one of BasicAuth, APIKey, or JWTToken is required")
	}
	if count > 1 {
		return nil, errors.New("only one of BasicAuth, APIKey, or JWTToken may be set")
	}

	switch {
	case cfg.BasicAuth != nil:
		if cfg.BasicAuth.Username == "" {
			return nil, errors.New("BasicAuth.Username is required")
		}
		return basicAuthEditor(cfg.BasicAuth.Username, cfg.BasicAuth.Password), nil
	case cfg.APIKey != "":
		return apiKeyEditor(cfg.APIKey), nil
	default:
		return bearerTokenEditor(cfg.JWTToken), nil
	}
}

func authMethodCount(cfg Config) int {
	n := 0
	if cfg.BasicAuth != nil {
		n++
	}
	if cfg.APIKey != "" {
		n++
	}
	if cfg.JWTToken != "" {
		n++
	}
	return n
}

func basicAuthEditor(username, password string) api.RequestEditorFn {
	return func(_ context.Context, req *http.Request) error {
		req.SetBasicAuth(username, password)
		return nil
	}
}

func apiKeyEditor(key string) api.RequestEditorFn {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", key))
		return nil
	}
}

func bearerTokenEditor(token string) api.RequestEditorFn {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}
}
