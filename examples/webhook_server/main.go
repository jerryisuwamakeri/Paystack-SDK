// Command webhook_server runs a minimal HTTP server that verifies and
// handles Paystack webhook events.
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	paystack "github.com/jerryisuwamakeri/Paystack-SDK"
)

func main() {
	secretKey := os.Getenv("PAYSTACK_SECRET_KEY")

	http.HandleFunc("/webhooks/paystack", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "could not read body", http.StatusBadRequest)
			return
		}

		signature := r.Header.Get("x-paystack-signature")
		event, err := paystack.VerifyAndParseWebhook(secretKey, signature, body)
		if err != nil {
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}

		switch event.Event {
		case paystack.EventChargeSuccess:
			log.Println("charge succeeded")
		case paystack.EventTransferSuccess:
			log.Println("transfer succeeded")
		case paystack.EventTransferFailed:
			log.Println("transfer failed")
		default:
			log.Println("unhandled event:", event.Event)
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
