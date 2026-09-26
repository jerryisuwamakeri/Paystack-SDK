package paystack

import (
	"net/http"
	"strings"
)

// redactedValue replaces sensitive values wherever this SDK produces
// diagnostic output.
const redactedValue = "[REDACTED]"

// sensitiveHeaders lists HTTP headers that must never appear in logs,
// error messages, or observability output in unredacted form. Keys are
// lower-cased for case-insensitive comparison.
var sensitiveHeaders = map[string]bool{
	"authorization":        true,
	"x-paystack-signature": true,
	"cookie":               true,
	"set-cookie":           true,
}

// redactHeaders returns a copy of h with sensitive header values replaced
// by a fixed placeholder, safe to include in logs, error messages, or
// observability hooks. The original header set is left untouched.
func redactHeaders(h http.Header) http.Header {
	redacted := make(http.Header, len(h))
	for key, values := range h {
		if sensitiveHeaders[strings.ToLower(key)] {
			redacted[key] = []string{redactedValue}
			continue
		}
		redacted[key] = append([]string(nil), values...)
	}
	return redacted
}
