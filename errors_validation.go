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

	ErrMissingName             = errors.New("paystack: name is required")
	ErrMissingAccountNumber    = errors.New("paystack: account number is required")
	ErrMissingBankCode         = errors.New("paystack: bank code is required")
	ErrMissingRecipientCode    = errors.New("paystack: recipient code is required")
	ErrMissingSource           = errors.New("paystack: source is required")
	ErrMissingReason           = errors.New("paystack: reason is required")
	ErrMissingTransferCode     = errors.New("paystack: transfer code is required")
	ErrMissingOTP              = errors.New("paystack: otp is required")
	ErrMissingInterval         = errors.New("paystack: interval is required")
	ErrMissingPlanCode         = errors.New("paystack: plan code is required")
	ErrMissingSubscriptionCode = errors.New("paystack: subscription code is required")
	ErrMissingEmailToken       = errors.New("paystack: email token is required")
	ErrMissingCustomer         = errors.New("paystack: customer is required")
)
