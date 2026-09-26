package paystack

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
)

// ErrInvalidWebhookSignature indicates that a webhook payload's signature
// did not match the signature computed from the configured secret, or that
// the supplied signature was empty or malformed.
var ErrInvalidWebhookSignature = errors.New("paystack: invalid webhook signature")

// VerifyWebhook validates that payload was sent by Paystack by recomputing
// its HMAC-SHA512 signature using secret (your Paystack secret key) and
// comparing it, in constant time, against signature (the value of the
// x-paystack-signature request header).
//
// Applications must verify every webhook before acting on its contents:
//
//	sig := r.Header.Get("x-paystack-signature")
//	body, _ := io.ReadAll(r.Body)
//	if err := paystack.VerifyWebhook(secretKey, sig, body); err != nil {
//	    http.Error(w, "invalid signature", http.StatusBadRequest)
//	    return
//	}
func VerifyWebhook(secret, signature string, payload []byte) error {
	if signature == "" {
		return ErrInvalidWebhookSignature
	}

	got, err := hex.DecodeString(signature)
	if err != nil {
		return ErrInvalidWebhookSignature
	}

	want, err := hex.DecodeString(computeWebhookSignature(secret, payload))
	if err != nil {
		return ErrInvalidWebhookSignature
	}

	if !hmac.Equal(got, want) {
		return ErrInvalidWebhookSignature
	}
	return nil
}

// SignWebhookPayload computes the HMAC-SHA512 signature Paystack would send
// in the x-paystack-signature header for payload, keyed with secret.
//
// This is primarily useful for testing webhook handlers; see the
// paystacktest package for higher-level test helpers built on top of it.
func SignWebhookPayload(secret string, payload []byte) string {
	return computeWebhookSignature(secret, payload)
}

// computeWebhookSignature returns the lowercase hex-encoded HMAC-SHA512 of
// payload, keyed with secret, matching Paystack's x-paystack-signature
// scheme.
func computeWebhookSignature(secret string, payload []byte) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
