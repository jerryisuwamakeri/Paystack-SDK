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

// transferData is the payload shape used by TransferSuccess and
// TransferFailed.
type transferData struct {
	Reference string `json:"reference"`
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
	Reason    string `json:"reason,omitempty"`
}

// TransferOption customizes a generated transfer test event.
type TransferOption func(*transferData)

// WithTransferReference overrides the transfer reference on a generated
// event.
func WithTransferReference(ref string) TransferOption {
	return func(d *transferData) { d.Reference = ref }
}

// WithTransferAmount overrides the amount, in the smallest currency unit,
// on a generated event.
func WithTransferAmount(amount int) TransferOption {
	return func(d *transferData) { d.Amount = amount }
}

// TransferSuccess returns a raw JSON payload shaped like a real
// transfer.success webhook event, for use in tests.
func TransferSuccess(opts ...TransferOption) []byte {
	data := transferData{
		Reference: "trf_test_000001",
		Amount:    250000,
		Currency:  "NGN",
		Status:    "success",
	}
	for _, opt := range opts {
		opt(&data)
	}
	return mustMarshalEvent("transfer.success", data)
}

// TransferFailed returns a raw JSON payload shaped like a real
// transfer.failed webhook event, for use in tests.
func TransferFailed(reason string, opts ...TransferOption) []byte {
	data := transferData{
		Reference: "trf_test_000001",
		Amount:    250000,
		Currency:  "NGN",
		Status:    "failed",
		Reason:    reason,
	}
	for _, opt := range opts {
		opt(&data)
	}
	return mustMarshalEvent("transfer.failed", data)
}

// refundData is the payload shape used by RefundProcessed.
type refundData struct {
	Reference      string `json:"reference"`
	TransactionRef string `json:"transaction_reference"`
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	Status         string `json:"status"`
}

// RefundOption customizes a generated refund test event.
type RefundOption func(*refundData)

// WithRefundAmount overrides the amount, in the smallest currency unit, on
// a generated event.
func WithRefundAmount(amount int) RefundOption {
	return func(d *refundData) { d.Amount = amount }
}

// RefundProcessed returns a raw JSON payload shaped like a real
// refund.processed webhook event, for use in tests.
func RefundProcessed(transactionReference string, opts ...RefundOption) []byte {
	data := refundData{
		Reference:      "ref_refund_000001",
		TransactionRef: transactionReference,
		Amount:         500000,
		Currency:       "NGN",
		Status:         "processed",
	}
	for _, opt := range opts {
		opt(&data)
	}
	return mustMarshalEvent("refund.processed", data)
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
