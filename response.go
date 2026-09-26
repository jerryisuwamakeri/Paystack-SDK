package paystack

// Response carries transport-level metadata about a completed API call,
// returned alongside the decoded result from every resource method.
type Response struct {
	// StatusCode is the HTTP status code of the final response.
	StatusCode int

	// RequestID is the Paystack-provided request identifier, when present.
	RequestID string

	// RateLimit carries rate-limit metadata from the response headers, when
	// the API supplied them.
	RateLimit RateLimit

	// Meta carries pagination metadata for list endpoints, when present.
	Meta *Meta
}

// Meta carries pagination metadata returned by Paystack list endpoints.
type Meta struct {
	Total     int    `json:"total,omitempty"`
	Skipped   int    `json:"skipped,omitempty"`
	PerPage   int    `json:"perPage,omitempty"`
	Page      int    `json:"page,omitempty"`
	PageCount int    `json:"pageCount,omitempty"`
	Next      string `json:"next,omitempty"`
	Previous  string `json:"previous,omitempty"`
}
