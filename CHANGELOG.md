# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/) once
v1.0.0 ships.

## [Unreleased]

### Added

- Core client with functional options: `WithBaseURL`, `WithTimeout`,
  `WithHTTPClient`, `WithRetryPolicy`, `WithObserver`.
- Structured `APIError` with sentinel error classification
  (`ErrUnauthorized`, `ErrRateLimited`, `ErrNotFound`, `ErrTimeout`).
- Automatic retries with exponential backoff and jitter for transient
  failures (429, 500, 502, 503, 504), respecting context cancellation.
- Per-request idempotency keys (`WithIdempotencyKey`) and request
  correlation IDs (`WithRequestID`).
- `TransportObserver` hooks for provider-agnostic observability.
- Rate-limit header parsing.
- Generic, loop-safe pagination `Iterator` with `ListAll` on every list
  endpoint.
- HMAC-SHA512 webhook signature verification with constant-time
  comparison (`VerifyWebhook`, `VerifyAndParseWebhook`) and typed
  `EventType` constants.
- `paystacktest` subpackage: webhook signing, canned event payloads
  (`ChargeSuccess`, `TransferSuccess`, `TransferFailed`,
  `RefundProcessed`), and a mock HTTP server for integration tests.
- Resource coverage: Transactions, Customers, Transfer Recipients,
  Transfers, Plans, Subscriptions, Refunds, Miscellaneous (banks and
  countries), Verification (account and BVN resolution), Products,
  Subaccounts, Settlements, Dedicated Virtual Accounts, Split Payments.
- Secret redaction for headers in logs and observability output; the
  `Client` type never leaks its secret key through `fmt` formatting.
- CI workflows for build/vet/test (with the race detector),
  vulnerability scanning, and secret detection.
