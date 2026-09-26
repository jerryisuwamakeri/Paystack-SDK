package paystack

import (
	"encoding/json"
	"fmt"
)

// Event is a parsed Paystack webhook payload. Data holds the event-specific
// payload as raw JSON so applications can decode it into a strongly typed
// struct for the events they care about, without this SDK needing to model
// every event shape up front.
type Event struct {
	Event EventType       `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// ParseWebhook decodes a webhook payload into an Event.
//
// ParseWebhook does not verify the payload's signature. Call VerifyWebhook
// first, or use VerifyAndParseWebhook, before trusting the event's
// contents.
func ParseWebhook(payload []byte) (*Event, error) {
	var evt Event
	if err := json.Unmarshal(payload, &evt); err != nil {
		return nil, fmt.Errorf("paystack: parsing webhook payload: %w", err)
	}
	return &evt, nil
}

// VerifyAndParseWebhook verifies payload's signature and, only once valid,
// parses it into an Event. This is the recommended entry point for webhook
// handlers, since it makes it impossible to accidentally act on an
// unverified payload.
func VerifyAndParseWebhook(secret, signature string, payload []byte) (*Event, error) {
	if err := VerifyWebhook(secret, signature, payload); err != nil {
		return nil, err
	}
	return ParseWebhook(payload)
}
