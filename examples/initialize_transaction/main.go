// Command initialize_transaction starts a hosted checkout transaction and
// prints the URL to redirect the customer to.
package main

import (
	"context"
	"log"
	"os"

	paystack "github.com/jerryisuwamakeri/Paystack-SDK"
)

func main() {
	client, err := paystack.NewClient(os.Getenv("PAYSTACK_SECRET_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	result, _, err := client.Transactions.Initialize(context.Background(), paystack.InitializeTransactionRequest{
		Email:  "customer@example.com",
		Amount: 500000, // NGN 5,000.00, in kobo
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("redirect the customer to:", result.AuthorizationURL)
	log.Println("reference:", result.Reference)
}
