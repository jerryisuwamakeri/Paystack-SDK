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
	ErrMissingTransaction      = errors.New("paystack: transaction is required")
	ErrMissingBVN              = errors.New("paystack: bvn is required")
	ErrMissingProductID        = errors.New("paystack: product id is required")
	ErrMissingBusinessName     = errors.New("paystack: business name is required")
	ErrMissingSettlementBank   = errors.New("paystack: settlement bank is required")
	ErrMissingSubaccountCode   = errors.New("paystack: subaccount code is required")
	ErrMissingSplitType        = errors.New("paystack: split type is required")
	ErrMissingSplitSubaccounts = errors.New("paystack: at least one subaccount is required")
	ErrMissingSplitCode        = errors.New("paystack: split code is required")
	ErrMissingBulkCharges      = errors.New("paystack: at least one charge is required")
	ErrMissingBatchCode        = errors.New("paystack: batch code is required")
	ErrMissingResolution       = errors.New("paystack: resolution is required")
	ErrMissingMessage          = errors.New("paystack: message is required")
	ErrMissingUploadedFilename = errors.New("paystack: uploaded filename is required")
	ErrMissingCustomerName     = errors.New("paystack: customer name is required")
	ErrMissingCustomerPhone    = errors.New("paystack: customer phone is required")
	ErrMissingServiceDetails   = errors.New("paystack: service details are required")
	ErrMissingDomainName       = errors.New("paystack: domain name is required")
	ErrMissingPIN              = errors.New("paystack: pin is required")
	ErrMissingPhone            = errors.New("paystack: phone is required")
	ErrMissingBirthday         = errors.New("paystack: birthday is required")
	ErrMissingTerminalID       = errors.New("paystack: terminal id is required")
	ErrMissingEventID          = errors.New("paystack: event id is required")
	ErrMissingEventType        = errors.New("paystack: event type is required")
	ErrMissingEventAction      = errors.New("paystack: event action is required")
)
