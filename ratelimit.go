package paystack

import (
	"net/http"
	"strconv"
	"time"
)

// RateLimit exposes rate-limit metadata parsed from response headers, when
// the API supplies them. A zero-value RateLimit means no rate-limit
// headers were present on the response.
type RateLimit struct {
	Limit     int
	Remaining int
	ResetAt   time.Time
}

func parseRateLimit(h http.Header) RateLimit {
	var rl RateLimit
	if v := h.Get("X-RateLimit-Limit"); v != "" {
		rl.Limit, _ = strconv.Atoi(v)
	}
	if v := h.Get("X-RateLimit-Remaining"); v != "" {
		rl.Remaining, _ = strconv.Atoi(v)
	}
	if v := h.Get("X-RateLimit-Reset"); v != "" {
		if secs, err := strconv.ParseInt(v, 10, 64); err == nil {
			rl.ResetAt = time.Unix(secs, 0)
		}
	}
	return rl
}
