package paystacktest

import (
	"encoding/json"
	"fmt"
	"time"
)

// chargeSuccessData is the payload shape used by ChargeSuccess.
type chargeSuccessData struct {
	Reference     string `json:"reference"`
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	PaidAt        string `json:"paid_at"`
	CustomerEmail string `json:"customer_email"`
}

// ChargeSuccessOption customizes a generated charge.success test event.
type ChargeSuccessOption func(*chargeSuccessData)

// WithReference overrides the transaction reference on a generated event.
func WithReference(ref string) ChargeSuccessOption {
	return func(d *chargeSuccessData) { d.Reference = ref }
}

// WithAmount overrides the amount, in the smallest currency unit, on a
// generated event.
func WithAmount(amount int) ChargeSuccessOption {
	return func(d *chargeSuccessData) { d.Amount = amount }
}

// WithCustomerEmail overrides the customer email on a generated event.
func WithCustomerEmail(email string) ChargeSuccessOption {
	return func(d *chargeSuccessData) { d.CustomerEmail = email }
}

// ChargeSuccess returns a raw JSON payload shaped like a real
// charge.success webhook event, for use in tests. Sensible defaults are
// used for any field not overridden with an option.
func ChargeSuccess(opts ...ChargeSuccessOption) []byte {
	data := chargeSuccessData{
		Reference:     "ref_test_000001",
		Amount:        500000,
		Currency:      "NGN",
		Status:        "success",
		PaidAt:        time.Now().UTC().Format(time.RFC3339),
		CustomerEmail: "customer@example.com",
	}
	for _, opt := range opts {
		opt(&data)
	}
	return mustMarshalEvent("charge.success", data)
}

func mustMarshalEvent(event string, data any) []byte {
	buf, err := json.Marshal(map[string]any{"event": event, "data": data})
	if err != nil {
		// Unreachable for the well-formed struct literals this package
		// constructs internally.
		panic(fmt.Sprintf("paystacktest: marshaling event: %v", err))
	}
	return buf
}
