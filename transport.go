package paystack

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"time"
)

// do executes an API call, applying the client's retry policy, invoking
// observability hooks, and decoding a successful response's data payload
// into out (which may be nil for calls that do not return a body).
func (c *Client) do(ctx context.Context, method, path string, body, out any, opts ...RequestOption) (*Response, error) {
	cfg := newRequestConfig(opts)
	maxAttempts := c.retryPolicy.MaxRetries + 1

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := c.newRequest(ctx, method, path, body, cfg)
		if err != nil {
			return nil, err
		}

		info := RequestInfo{
			Method:   method,
			Endpoint: path,
			Headers:  redactHeaders(req.Header),
			Attempt:  attempt,
		}
		c.observer.RequestStarted(info)

		start := time.Now()
		resp, err := c.httpClient.Do(req)
		duration := time.Since(start)

		if err != nil {
			c.observer.RequestFailed(info, err)

			if ctx.Err() != nil {
				return nil, &APIError{Endpoint: path, Message: "request canceled or timed out", Err: ctx.Err()}
			}

			apiErr := &APIError{Endpoint: path, Message: err.Error(), Err: err, Retryable: isRetryableNetErr(err)}
			lastErr = apiErr

			if apiErr.Retryable && attempt < maxAttempts {
				if werr := c.sleepForRetry(ctx, attempt, info, apiErr); werr != nil {
					return nil, &APIError{Endpoint: path, Message: "context canceled during retry backoff", Err: werr}
				}
				continue
			}
			return nil, apiErr
		}

		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		requestID := resp.Header.Get(headerRequestID)

		c.observer.RequestCompleted(ResponseInfo{
			RequestInfo: info,
			StatusCode:  resp.StatusCode,
			Duration:    duration,
			RequestID:   requestID,
		})

		if readErr != nil {
			return nil, &APIError{StatusCode: resp.StatusCode, Endpoint: path, Message: "reading response body", Err: readErr, RequestID: requestID}
		}

		result := &Response{
			StatusCode: resp.StatusCode,
			RequestID:  requestID,
			RateLimit:  parseRateLimit(resp.Header),
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if decodeErr := decodeInto(respBody, out, result); decodeErr != nil {
				return result, &APIError{StatusCode: resp.StatusCode, Endpoint: path, Message: "decoding response", Err: decodeErr, RequestID: requestID}
			}
			return result, nil
		}

		apiErr := buildAPIError(resp.StatusCode, path, requestID, respBody)
		lastErr = apiErr

		if apiErr.Retryable && attempt < maxAttempts {
			if werr := c.sleepForRetry(ctx, attempt, info, apiErr); werr != nil {
				return result, &APIError{Endpoint: path, Message: "context canceled during retry backoff", Err: werr, RequestID: requestID}
			}
			continue
		}

		return result, apiErr
	}

	return nil, lastErr
}

// decodeInto unmarshals a successful response body's envelope, populating
// out with the "data" field and result.Meta with any pagination metadata.
func decodeInto(body []byte, out any, result *Response) error {
	if len(body) == 0 {
		return nil
	}

	var envelope apiEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return err
	}
	result.Meta = envelope.Meta

	if out == nil || len(envelope.Data) == 0 {
		return nil
	}
	return json.Unmarshal(envelope.Data, out)
}

// sleepForRetry notifies the observer of a scheduled retry and blocks for
// the backoff delay, returning early with ctx.Err() if the context is
// canceled before the delay elapses.
func (c *Client) sleepForRetry(ctx context.Context, attempt int, info RequestInfo, cause error) error {
	delay := c.retryPolicy.Backoff.Delay(attempt)
	c.observer.Retry(RetryInfo{RequestInfo: info, Delay: delay, Cause: cause})

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// isRetryableNetErr classifies a transport-level error (as opposed to an
// HTTP status code) as safe to retry. Context cancellation and deadlines
// are never retried here since they represent an explicit caller decision.
func isRetryableNetErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	// Unknown transport failures (connection refused, DNS failure,
	// unexpected EOF) are generally transient in production networks.
	return true
}
