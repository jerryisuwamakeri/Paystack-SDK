package paystack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	headerAuthorization  = "Authorization"
	headerIdempotencyKey = "Idempotency-Key"
	headerRequestID      = "X-Request-ID"
	headerContentType    = "Content-Type"
	headerUserAgent      = "User-Agent"
)

// newRequest builds an *http.Request for the given method and path,
// encoding body as JSON when non-nil, and attaching authentication and
// per-request headers.
func (c *Client) newRequest(ctx context.Context, method, path string, body any, cfg requestConfig) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("paystack: encoding request body: %w", err)
		}
		reader = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, fmt.Errorf("paystack: building request: %w", err)
	}

	req.Header.Set(headerAuthorization, "Bearer "+c.secretKey)
	req.Header.Set(headerUserAgent, userAgent)
	if body != nil {
		req.Header.Set(headerContentType, "application/json")
	}
	if cfg.idempotencyKey != "" {
		req.Header.Set(headerIdempotencyKey, cfg.idempotencyKey)
	}
	if cfg.requestID != "" {
		req.Header.Set(headerRequestID, cfg.requestID)
	}

	return req, nil
}

// apiEnvelope mirrors the standard {status, message, data, meta} envelope
// that Paystack wraps every response in.
type apiEnvelope struct {
	Status  bool            `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Meta    *Meta           `json:"meta,omitempty"`
}

// buildAPIError constructs an APIError from a non-2xx HTTP response body,
// falling back to the standard HTTP status text if the body cannot be
// parsed as the standard error envelope.
func buildAPIError(statusCode int, endpoint, requestID string, body []byte) *APIError {
	message := http.StatusText(statusCode)

	var envelope apiEnvelope
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Message != "" {
		message = envelope.Message
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Endpoint:   endpoint,
		RequestID:  requestID,
		Retryable:  isRetryableStatus(statusCode),
	}
}
