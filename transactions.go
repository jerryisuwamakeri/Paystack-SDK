package paystack

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// TransactionService groups the transaction-related Paystack API
// operations. Access it via Client.Transactions.
type TransactionService struct {
	client *Client
}

// TransactionCustomer is the customer information embedded in a
// Transaction.
type TransactionCustomer struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	CustomerCode string `json:"customer_code"`
}

// Transaction represents a Paystack transaction record.
type Transaction struct {
	ID              int64               `json:"id"`
	Reference       string              `json:"reference"`
	Amount          int64               `json:"amount"`
	Currency        string              `json:"currency"`
	Status          string              `json:"status"`
	GatewayResponse string              `json:"gateway_response"`
	Channel         string              `json:"channel"`
	PaidAt          *time.Time          `json:"paid_at,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	Metadata        map[string]any      `json:"metadata,omitempty"`
	Customer        TransactionCustomer `json:"customer"`
}

// InitializeTransactionRequest is the payload for
// TransactionService.Initialize.
type InitializeTransactionRequest struct {
	// Email is the customer's email address. Required.
	Email string `json:"email"`

	// Amount is the amount to charge, in the smallest currency unit (for
	// example, kobo for NGN or cents for USD). Required.
	Amount int64 `json:"amount"`

	// Currency is the three-letter ISO currency code. Optional; defaults to
	// the integration's currency.
	Currency string `json:"currency,omitempty"`

	// Reference is a unique transaction reference. Optional; Paystack
	// generates one when omitted.
	Reference string `json:"reference,omitempty"`

	// CallbackURL overrides the integration's default callback URL for this
	// transaction. Optional.
	CallbackURL string `json:"callback_url,omitempty"`

	// Metadata carries arbitrary structured data through to the
	// transaction record. Optional.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// InitializeTransactionResponse is returned by
// TransactionService.Initialize.
type InitializeTransactionResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}

// Initialize starts a new transaction, returning a checkout URL the
// customer should be redirected to.
func (s *TransactionService) Initialize(ctx context.Context, req InitializeTransactionRequest, opts ...RequestOption) (*InitializeTransactionResponse, *Response, error) {
	if req.Email == "" {
		return nil, nil, ErrMissingEmail
	}
	if req.Amount <= 0 {
		return nil, nil, ErrInvalidAmount
	}

	var out InitializeTransactionResponse
	resp, err := s.client.do(ctx, "POST", "/transaction/initialize", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Verify retrieves a transaction's current status by its reference. This
// is the source of truth for whether a payment succeeded and should be
// called from the server side rather than trusting the client-side
// redirect alone.
func (s *TransactionService) Verify(ctx context.Context, reference string, opts ...RequestOption) (*Transaction, *Response, error) {
	if reference == "" {
		return nil, nil, ErrMissingReference
	}

	var out Transaction
	resp, err := s.client.do(ctx, "GET", "/transaction/verify/"+url.PathEscape(reference), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a single transaction by its numeric ID.
func (s *TransactionService) Fetch(ctx context.Context, id int64, opts ...RequestOption) (*Transaction, *Response, error) {
	var out Transaction
	resp, err := s.client.do(ctx, "GET", fmt.Sprintf("/transaction/%d", id), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ListTransactionsParams filters TransactionService.List and
// TransactionService.ListAll.
type ListTransactionsParams struct {
	Page     int
	PerPage  int
	Customer int64
	Status   string
	From     time.Time
	To       time.Time
}

func (p ListTransactionsParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	if p.Customer > 0 {
		q.Set("customer", strconv.FormatInt(p.Customer, 10))
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	if !p.From.IsZero() {
		q.Set("from", p.From.Format(time.RFC3339))
	}
	if !p.To.IsZero() {
		q.Set("to", p.To.Format(time.RFC3339))
	}
	return q
}

// List retrieves a single page of transactions matching params.
func (s *TransactionService) List(ctx context.Context, params ListTransactionsParams, opts ...RequestOption) ([]Transaction, *Response, error) {
	path := "/transaction"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Transaction
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every transaction matching params,
// fetching additional pages on demand as the caller advances it.
func (s *TransactionService) ListAll(ctx context.Context, params ListTransactionsParams, opts ...RequestOption) *Iterator[Transaction] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Transaction, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
