package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// TransferService groups the transfer-related Paystack API operations.
// Access it via Client.Transfers.
//
// Transfer creation supports idempotency keys, since resubmitting a
// transfer request after a lost response could otherwise move money
// twice:
//
//	transfer, _, err := client.Transfers.Create(ctx, req,
//	    paystack.WithIdempotencyKey("transfer-2026-000001"),
//	)
type TransferService struct {
	client *Client
}

// Transfer represents a Paystack transfer record.
type Transfer struct {
	ID           int64              `json:"id"`
	TransferCode string             `json:"transfer_code"`
	Reference    string             `json:"reference"`
	Amount       int64              `json:"amount"`
	Currency     string             `json:"currency"`
	Status       string             `json:"status"`
	Reason       string             `json:"reason"`
	Recipient    *TransferRecipient `json:"recipient,omitempty"`
}

// CreateTransferRequest is the payload for TransferService.Create.
type CreateTransferRequest struct {
	// Source is the origin of the funds. Defaults to "balance" when empty.
	Source string `json:"source"`

	// Amount is the amount to transfer, in the smallest currency unit.
	// Required.
	Amount int64 `json:"amount"`

	// Recipient is the recipient code of a previously created
	// TransferRecipient. Required.
	Recipient string `json:"recipient"`

	Reason    string `json:"reason,omitempty"`
	Currency  string `json:"currency,omitempty"`
	Reference string `json:"reference,omitempty"`
}

// Create initiates a transfer to a saved recipient. Callers should supply
// WithIdempotencyKey to make retrying a Create call safe.
func (s *TransferService) Create(ctx context.Context, req CreateTransferRequest, opts ...RequestOption) (*Transfer, *Response, error) {
	if req.Amount <= 0 {
		return nil, nil, ErrInvalidAmount
	}
	if req.Recipient == "" {
		return nil, nil, ErrMissingRecipientCode
	}
	if req.Source == "" {
		req.Source = "balance"
	}

	var out Transfer
	resp, err := s.client.do(ctx, "POST", "/transfer", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// FinalizeTransferRequest is the payload for TransferService.Finalize.
type FinalizeTransferRequest struct {
	TransferCode string `json:"transfer_code"`
	OTP          string `json:"otp"`
}

// Finalize completes a transfer that requires OTP confirmation.
func (s *TransferService) Finalize(ctx context.Context, req FinalizeTransferRequest, opts ...RequestOption) (*Transfer, *Response, error) {
	if req.TransferCode == "" {
		return nil, nil, ErrMissingTransferCode
	}
	if req.OTP == "" {
		return nil, nil, ErrMissingOTP
	}

	var out Transfer
	resp, err := s.client.do(ctx, "POST", "/transfer/finalize_transfer", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a transfer by its numeric ID or transfer code.
func (s *TransferService) Fetch(ctx context.Context, idOrCode string, opts ...RequestOption) (*Transfer, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingTransferCode
	}

	var out Transfer
	resp, err := s.client.do(ctx, "GET", "/transfer/"+url.PathEscape(idOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ListTransfersParams filters TransferService.List and
// TransferService.ListAll.
type ListTransfersParams struct {
	Page    int
	PerPage int
	Status  string
}

func (p ListTransfersParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	if p.Status != "" {
		q.Set("status", p.Status)
	}
	return q
}

// List retrieves a single page of transfers.
func (s *TransferService) List(ctx context.Context, params ListTransfersParams, opts ...RequestOption) ([]Transfer, *Response, error) {
	path := "/transfer"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Transfer
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every transfer matching params, fetching
// additional pages on demand as the caller advances it.
func (s *TransferService) ListAll(ctx context.Context, params ListTransfersParams, opts ...RequestOption) *Iterator[Transfer] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Transfer, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
