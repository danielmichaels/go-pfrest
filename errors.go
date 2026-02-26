package pfrest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIError struct {
	Code       int    `json:"code"`
	Status     string `json:"status"`
	ResponseID string `json:"response_id"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	if e.ResponseID != "" {
		return fmt.Sprintf("pfrest: %d %s [%s]: %s", e.Code, e.Status, e.ResponseID, e.Message)
	}
	return fmt.Sprintf("pfrest: %d %s: %s", e.Code, e.Status, e.Message)
}

func CheckResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			Code:    resp.StatusCode,
			Status:  resp.Status,
			Message: "failed to read response body",
		}
	}

	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return &APIError{
			Code:    resp.StatusCode,
			Status:  resp.Status,
			Message: string(body),
		}
	}

	if apiErr.Code == 0 {
		apiErr.Code = resp.StatusCode
	}
	if apiErr.Status == "" {
		apiErr.Status = resp.Status
	}

	return &apiErr
}
