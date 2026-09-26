package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// RefundService groups the refund-related Paystack API operations. Access
// it via Client.Refunds.
type RefundService struct {
	client *Client
}

// Refund represents a Paystack refund record.
type Refund struct {
	ID                   int64  `json:"id"`
	TransactionReference string `json:"transaction_reference,omitempty"`
	Amount               int64  `json:"amount"`
	Currency             string `json:"currency"`
	Status               string `json:"status"`
}

// CreateRefundRequest is the payload for RefundService.Create.
type CreateRefundRequest struct {
	// Transaction is the transaction ID or reference to refund. Required.
	Transaction string `json:"transaction"`

	// Amount is the amount to refund, in the smallest currency unit.
	// Optional; omit for a full refund.
	Amount int64 `json:"amount,omitempty"`

	Currency     string `json:"currency,omitempty"`
	CustomerNote string `json:"customer_note,omitempty"`
	MerchantNote string `json:"merchant_note,omitempty"`
}

// Create issues a refund for a transaction, in full or in part.
func (s *RefundService) Create(ctx context.Context, req CreateRefundRequest, opts ...RequestOption) (*Refund, *Response, error) {
	if req.Transaction == "" {
		return nil, nil, ErrMissingTransaction
	}

	var out Refund
	resp, err := s.client.do(ctx, "POST", "/refund", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a refund by its numeric ID.
func (s *RefundService) Fetch(ctx context.Context, id string, opts ...RequestOption) (*Refund, *Response, error) {
	if id == "" {
		return nil, nil, ErrMissingReference
	}

	var out Refund
	resp, err := s.client.do(ctx, "GET", "/refund/"+url.PathEscape(id), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ListRefundsParams filters RefundService.List and RefundService.ListAll.
type ListRefundsParams struct {
	Page      int
	PerPage   int
	Reference string
}

func (p ListRefundsParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	if p.Reference != "" {
		q.Set("reference", p.Reference)
	}
	return q
}

// List retrieves a single page of refunds.
func (s *RefundService) List(ctx context.Context, params ListRefundsParams, opts ...RequestOption) ([]Refund, *Response, error) {
	path := "/refund"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Refund
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every refund matching params, fetching
// additional pages on demand as the caller advances it.
func (s *RefundService) ListAll(ctx context.Context, params ListRefundsParams, opts ...RequestOption) *Iterator[Refund] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Refund, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
