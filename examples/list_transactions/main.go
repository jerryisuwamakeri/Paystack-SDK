// Command list_transactions paginates through every successful transaction
// using the generic pagination Iterator.
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

	ctx := context.Background()
	iter := client.Transactions.ListAll(ctx, paystack.ListTransactionsParams{
		PerPage: 50,
		Status:  "success",
	})

	count := 0
	for iter.Next() {
		tx := iter.Item()
		log.Printf("%s: %d %s", tx.Reference, tx.Amount, tx.Currency)
		count++
	}
	if err := iter.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("processed %d transactions", count)
}
