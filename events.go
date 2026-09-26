package paystack

// EventType identifies the kind of webhook event sent by Paystack. New
// event types may be introduced by Paystack at any time; applications
// should handle an unrecognized EventType gracefully rather than treating
// it as an error.
type EventType string

const (
	EventChargeSuccess EventType = "charge.success"

	EventTransferSuccess  EventType = "transfer.success"
	EventTransferFailed   EventType = "transfer.failed"
	EventTransferReversed EventType = "transfer.reversed"

	EventRefundProcessed EventType = "refund.processed"
	EventRefundFailed    EventType = "refund.failed"

	EventSubscriptionCreate        EventType = "subscription.create"
	EventSubscriptionDisable       EventType = "subscription.disable"
	EventSubscriptionNotRenew      EventType = "subscription.not_renew"
	EventSubscriptionExpiringCards EventType = "subscription.expiring_cards"

	EventInvoiceCreate        EventType = "invoice.create"
	EventInvoiceUpdate        EventType = "invoice.update"
	EventInvoicePaymentFailed EventType = "invoice.payment_failed"

	EventCustomerIdentificationSuccess EventType = "customeridentification.success"
	EventCustomerIdentificationFailed  EventType = "customeridentification.failed"

	EventDedicatedAccountAssignSuccess EventType = "dedicatedaccount.assign.success"
	EventDedicatedAccountAssignFailed  EventType = "dedicatedaccount.assign.failed"
)
