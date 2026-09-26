# Paystack Go SDK

[![Test](https://github.com/jerryisuwamakeri/Paystack-SDK/actions/workflows/test.yml/badge.svg)](https://github.com/jerryisuwamakeri/Paystack-SDK/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jerryisuwamakeri/Paystack-SDK.svg)](https://pkg.go.dev/github.com/jerryisuwamakeri/Paystack-SDK)

An enterprise-grade Go client for the [Paystack API](https://paystack.com/docs/api/), built as a thin, reliable integration layer: it handles authentication, retries, timeouts, pagination, webhook verification, and error classification, while leaving business logic, database state, and observability providers entirely under your control.

## Installation

```bash
go get github.com/jerryisuwamakeri/Paystack-SDK
```

## Quick start

```go
package main

import (
	"context"
	"log"
	"os"

	paystack "github.com/jerryisuwamakeri/Paystack-SDK"
)

func main() {
	client, err := paystack.NewClient(os.Getenv("PAYSTACK_SECRET_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	result, _, err := client.Transactions.Initialize(ctx, paystack.InitializeTransactionRequest{
		Email:  "customer@example.com",
		Amount: 500000, // NGN 5,000.00, in kobo
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("redirect the customer to:", result.AuthorizationURL)
}
```

The secret key is read from runtime configuration; the SDK never persists, logs, or prints it (even `fmt.Printf("%+v", client)` is safe).

## Configuration

```go
client, err := paystack.NewClient(
	os.Getenv("PAYSTACK_SECRET_KEY"),
	paystack.WithBaseURL("https://api.paystack.co"), // override for a sandbox or test double
	paystack.WithTimeout(15*time.Second),
	paystack.WithRetryPolicy(paystack.RetryPolicy{
		MaxRetries: 3,
		Backoff:    paystack.ExponentialBackoff{Base: 250 * time.Millisecond, Max: 5 * time.Second},
	}),
	paystack.WithObserver(myObserver), // OpenTelemetry, Prometheus, Datadog, logging, ...
	paystack.WithHTTPClient(myHTTPClient), // custom transport, proxy, connection pooling
)
```

A `*Client` is immutable after construction and safe for concurrent use across your entire application.

## Idempotent transfers

Financial mutations accept per-call idempotency keys, so a retried request after a lost response never moves money twice:

```go
transfer, _, err := client.Transfers.Create(ctx,
	paystack.CreateTransferRequest{
		Amount:    500000,
		Recipient: "RCP_1a2b3c",
	},
	paystack.WithIdempotencyKey("transfer-2026-000001"),
)
```

Retries triggered internally by the SDK's retry policy always reuse the same key you provided.

## Pagination

`ListAll` returns a generic `Iterator` that fetches pages on demand, stops on an empty page or a caller-supplied page limit, and detects an API that repeats the same page instead of looping forever:

```go
iter := client.Transactions.ListAll(ctx, paystack.ListTransactionsParams{PerPage: 50})
for iter.Next() {
	tx := iter.Item()
	// process tx
}
if err := iter.Err(); err != nil {
	log.Fatal(err)
}
```

## Webhooks

```go
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	sig := r.Header.Get("x-paystack-signature")

	event, err := paystack.VerifyAndParseWebhook(secretKey, sig, body)
	if err != nil {
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	switch event.Event {
	case paystack.EventChargeSuccess:
		// decode event.Data and update your own records
	}

	w.WriteHeader(http.StatusOK)
}
```

`VerifyAndParseWebhook` uses HMAC-SHA512 with a constant-time comparison and rejects empty or malformed signatures before your handler ever sees the payload.

## Testing your integration

The `paystacktest` subpackage provides webhook signing, canned event payloads, and a mock HTTP server so you can test against this SDK without contacting the real API:

```go
import "github.com/jerryisuwamakeri/Paystack-SDK/paystacktest"

payload := paystacktest.ChargeSuccess(paystacktest.WithReference("ref_123"))
signature := paystacktest.Sign(payload, secretKey)
```

## Error handling

```go
var apiErr *paystack.APIError
if errors.As(err, &apiErr) {
	if apiErr.Retryable {
		// safe to retry
	}
	log.Printf("paystack error: %s (status %d, request %s)", apiErr.Message, apiErr.StatusCode, apiErr.RequestID)
}

if errors.Is(err, paystack.ErrRateLimited) {
	// back off
}
```

## Compatibility

| SDK  | Paystack API | Go    |
| ---- | ------------ | ----- |
| v0.x | Current      | 1.22+ |

Once v1.0.0 ships, existing code will continue to compile across minor releases; breaking changes will require a major version bump.

## Documentation

- [Paystack API reference](https://paystack.com/docs/api/)
- [Paystack docs](https://paystack.com/docs/)
- [Go package reference](https://pkg.go.dev/github.com/jerryisuwamakeri/Paystack-SDK)

## License

[MIT](LICENSE)
