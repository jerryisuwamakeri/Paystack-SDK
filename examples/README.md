# Examples

Each subdirectory is a standalone, runnable program. Set `PAYSTACK_SECRET_KEY`
before running any of them:

```bash
export PAYSTACK_SECRET_KEY=sk_test_xxxxx
go run ./examples/initialize_transaction
go run ./examples/verify_transaction <reference>
go run ./examples/list_transactions
go run ./examples/create_transfer <recipient_code> <idempotency_key>
go run ./examples/webhook_server
```

| Example                 | Demonstrates                                  |
| ------------------------ | ---------------------------------------------- |
| `initialize_transaction` | Starting a hosted checkout transaction         |
| `verify_transaction`     | Server-side verification of a completed charge |
| `list_transactions`      | The generic pagination `Iterator`              |
| `create_transfer`        | Idempotent transfer creation                   |
| `webhook_server`         | Webhook signature verification and event types |
