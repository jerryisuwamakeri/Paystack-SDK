package paystack

import "errors"

// Client-side validation sentinel errors. The SDK performs only cheap,
// obvious checks before making a request so that mistakes are caught
// without a round trip; the Paystack API remains the authoritative source
// for business-rule validation.
var (
	ErrMissingEmail     = errors.New("paystack: email is required")
	ErrInvalidAmount    = errors.New("paystack: amount must be greater than zero")
	ErrMissingReference = errors.New("paystack: reference is required")
)
