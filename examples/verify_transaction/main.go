// Command verify_transaction verifies a transaction by its reference and
// prints its status. Run this from your server after a customer completes
// checkout; never trust the client-side redirect alone.
package main

import (
	"context"
	"log"
	"os"

	paystack "github.com/jerryisuwamakeri/Paystack-SDK"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: verify_transaction <reference>")
	}
	reference := os.Args[1]

	client, err := paystack.NewClient(os.Getenv("PAYSTACK_SECRET_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	tx, _, err := client.Transactions.Verify(context.Background(), reference)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%s amount=%d currency=%s", tx.Status, tx.Amount, tx.Currency)
}
