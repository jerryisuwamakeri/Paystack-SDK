# Contributing

Thanks for considering a contribution to the Paystack Go SDK.

## Development setup

Requires Go 1.22 or later.

```bash
git clone https://github.com/jerryisuwamakeri/Paystack-SDK.git
cd Paystack-SDK
go build ./...
go test ./...
```

## Before opening a pull request

- `go build ./...` and `go vet ./...` must pass.
- `gofmt -l .` must print nothing (run `gofmt -w .` to fix formatting).
- `go test ./... -race` must pass.
- New behavior needs test coverage; bug fixes should include a test that
  fails before the fix and passes after it.

## Design principles

- **Make the safe path the easy path.** Retries, backoff, pagination,
  webhook verification, and secret redaction should be automatic, not
  something callers have to reimplement.
- **Thin integration layer, not a framework.** The SDK owns HTTP transport,
  request/response shape, and error classification. Business logic,
  database state, and retry *policy* decisions stay with the caller.
- **Minimal dependencies.** Prefer the standard library. A new external
  dependency needs a clear justification in the pull request description.
- **Monetary values are integers.** Amounts are always the smallest
  currency unit (e.g. kobo, cents) as `int64`, never `float64`.

## Commit messages

Describe the change and, where it isn't obvious, why it was made. Small,
focused commits are preferred over large mixed ones.

## Reporting security issues

Please see [SECURITY.md](SECURITY.md) instead of opening a public issue.
