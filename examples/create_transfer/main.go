// Command create_transfer creates a transfer to a saved recipient using an
// idempotency key, so retrying the program after a lost response never
// moves money twice.
package main

import (
	"context"
	"log"
	"os"

	paystack "github.com/jerryisuwamakeri/Paystack-SDK"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: create_transfer <recipient_code> <idempotency_key>")
	}
	recipientCode := os.Args[1]
	idempotencyKey := os.Args[2]

	client, err := paystack.NewClient(os.Getenv("PAYSTACK_SECRET_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	transfer, _, err := client.Transfers.Create(context.Background(),
		paystack.CreateTransferRequest{
			Amount:    500000,
			Recipient: recipientCode,
			Reason:    "Example payout",
		},
		paystack.WithIdempotencyKey(idempotencyKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("transfer_code=%s status=%s", transfer.TransferCode, transfer.Status)
}
