package pfrest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, wantAuth string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("Authorization")
		if wantAuth != "" && got != wantAuth {
			t.Errorf("Authorization header = %q, want %q", got, wantAuth)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    200,
			"status":  "ok",
			"message": "test",
			"data":    map[string]any{},
		})
	}))
}

func TestNewClient_MissingBaseURL(t *testing.T) {
	_, err := NewClient(Config{APIKey: "test"})
	if err == nil {
		t.Fatal("expected error for missing BaseURL")
	}
}

func TestNewClient_NoAuth(t *testing.T) {
	_, err := NewClient(Config{BaseURL: "http://localhost"})
	if err == nil {
		t.Fatal("expected error when no auth method set")
	}
}

func TestNewClient_ConflictingAuth(t *testing.T) {
	_, err := NewClient(Config{
		BaseURL:  "http://localhost",
		APIKey:   "key",
		JWTToken: "token",
	})
	if err == nil {
		t.Fatal("expected error for conflicting auth methods")
	}
}

func TestNewClient_BasicAuth(t *testing.T) {
	ts := newTestServer(t, "")
	defer ts.Close()

	client, err := NewClient(Config{
		BaseURL: ts.URL,
		BasicAuth: &BasicAuthConfig{
			Username: "admin",
			Password: "pfsense",
		},
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client.Raw() == nil {
		t.Fatal("Raw() returned nil")
	}
}

func TestNewClient_BasicAuth_MissingUsername(t *testing.T) {
	_, err := NewClient(Config{
		BaseURL:   "http://localhost",
		BasicAuth: &BasicAuthConfig{Password: "pass"},
	})
	if err == nil {
		t.Fatal("expected error for missing BasicAuth.Username")
	}
}

func TestNewClient_APIKey(t *testing.T) {
	ts := newTestServer(t, "Token my-api-key")
	defer ts.Close()

	client, err := NewClient(Config{
		BaseURL: ts.URL,
		APIKey:  "my-api-key",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client.Raw() == nil {
		t.Fatal("Raw() returned nil")
	}
}

func TestNewClient_JWTToken(t *testing.T) {
	ts := newTestServer(t, "Bearer my-jwt-token")
	defer ts.Close()

	client, err := NewClient(Config{
		BaseURL:  ts.URL,
		JWTToken: "my-jwt-token",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client.Raw() == nil {
		t.Fatal("Raw() returned nil")
	}
}

func TestNewClient_TrailingSlash(t *testing.T) {
	ts := newTestServer(t, "")
	defer ts.Close()

	client, err := NewClient(Config{
		BaseURL: ts.URL + "/",
		APIKey:  "test",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client.Raw() == nil {
		t.Fatal("Raw() returned nil")
	}
}

func TestNewClient_CustomHTTPClient(t *testing.T) {
	ts := newTestServer(t, "")
	defer ts.Close()

	custom := &http.Client{}
	client, err := NewClient(Config{
		BaseURL:    ts.URL,
		APIKey:     "test",
		HTTPClient: custom,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if client.Raw() == nil {
		t.Fatal("Raw() returned nil")
	}
}
