// Package paystacktest provides test infrastructure for applications that
// integrate with Paystack: webhook signing helpers, canned event payloads,
// and a mock HTTP server, so that webhook handlers and API integrations can
// be tested without contacting the real Paystack API.
package paystacktest

import "github.com/jerryisuwamakeri/Paystack-SDK"

// Sign computes the x-paystack-signature value Paystack would send for a
// webhook payload signed with secret, for use in tests:
//
//	payload := paystacktest.ChargeSuccess()
//	signature := paystacktest.Sign(payload, secret)
func Sign(payload []byte, secret string) string {
	return paystack.SignWebhookPayload(secret, payload)
}
