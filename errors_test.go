package pfrest

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "with response_id",
			err:  APIError{Code: 400, Status: "bad request", ResponseID: "VALIDATION_ERROR", Message: "invalid input"},
			want: "pfrest: 400 bad request [VALIDATION_ERROR]: invalid input",
		},
		{
			name: "without response_id",
			err:  APIError{Code: 500, Status: "internal server error", Message: "something broke"},
			want: "pfrest: 500 internal server error: something broke",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCheckResponse_Success(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	if err := CheckResponse(resp); err != nil {
		t.Fatalf("CheckResponse(200) = %v, want nil", err)
	}
}

func TestCheckResponse_ErrorEnvelope(t *testing.T) {
	body := `{"code":400,"status":"bad request","response_id":"VALIDATION_ERROR","message":"field x is required"}`
	resp := &http.Response{
		StatusCode: 400,
		Status:     "400 Bad Request",
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("CheckResponse(400) = nil, want error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.Code != 400 {
		t.Errorf("Code = %d, want 400", apiErr.Code)
	}
	if apiErr.ResponseID != "VALIDATION_ERROR" {
		t.Errorf("ResponseID = %q, want VALIDATION_ERROR", apiErr.ResponseID)
	}
	if apiErr.Message != "field x is required" {
		t.Errorf("Message = %q, want 'field x is required'", apiErr.Message)
	}
}

func TestCheckResponse_NonJSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: 502,
		Status:     "502 Bad Gateway",
		Body:       io.NopCloser(strings.NewReader("<html>Bad Gateway</html>")),
	}

	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("CheckResponse(502) = nil, want error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.Code != 502 {
		t.Errorf("Code = %d, want 502", apiErr.Code)
	}
	if apiErr.Message != "<html>Bad Gateway</html>" {
		t.Errorf("Message = %q, want raw HTML body", apiErr.Message)
	}
}

func TestCheckResponse_401(t *testing.T) {
	body := `{"code":401,"status":"unauthorized","response_id":"AUTH_INVALID_CREDS","message":"invalid credentials"}`
	resp := &http.Response{
		StatusCode: 401,
		Status:     "401 Unauthorized",
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.ResponseID != "AUTH_INVALID_CREDS" {
		t.Errorf("ResponseID = %q, want AUTH_INVALID_CREDS", apiErr.ResponseID)
	}
}
