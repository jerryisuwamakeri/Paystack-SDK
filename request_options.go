package paystack

// RequestOption configures a single API call, as opposed to Option, which
// configures a Client for its entire lifetime.
type RequestOption func(*requestConfig)

type requestConfig struct {
	idempotencyKey string
	requestID      string
}

func newRequestConfig(opts []RequestOption) requestConfig {
	var cfg requestConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// WithIdempotencyKey attaches an idempotency key to a single request so
// that Paystack can safely deduplicate retried or resubmitted calls for
// operations that support idempotency, such as transfer creation.
//
// The same key must be reused across retries of the same logical
// operation:
//
//	key := "transfer-2026-000001"
//	transfer, err := client.Transfers.Create(ctx, req, paystack.WithIdempotencyKey(key))
//
// The SDK never generates a key on the caller's behalf and never mutates a
// caller-provided key, including during automatic retries.
func WithIdempotencyKey(key string) RequestOption {
	return func(cfg *requestConfig) {
		cfg.idempotencyKey = key
	}
}

// WithRequestID attaches an application-level correlation identifier to a
// single request. The identifier is echoed back on any resulting APIError
// and passed to observability hooks, so that application logs, SDK
// requests, and Paystack requests can be correlated during incident
// investigation.
func WithRequestID(id string) RequestOption {
	return func(cfg *requestConfig) {
		cfg.requestID = id
	}
}
