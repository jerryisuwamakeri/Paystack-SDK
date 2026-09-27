package paystack

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// InvoiceService groups the payment-request (invoice) API operations.
// Access it via Client.Invoices.
type InvoiceService struct {
	client *Client
}

// Invoice represents a Paystack payment request.
type Invoice struct {
	ID          int64      `json:"id"`
	RequestCode string     `json:"request_code"`
	Description string     `json:"description,omitempty"`
	Amount      int64      `json:"amount"`
	Currency    string     `json:"currency"`
	Status      string     `json:"status"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// CreateInvoiceRequest is the payload for InvoiceService.Create.
type CreateInvoiceRequest struct {
	// Customer is the customer's email address or customer code. Required.
	Customer string `json:"customer"`

	// Amount is optional; omit it for a payment request with a
	// customer-entered amount.
	Amount           int64  `json:"amount,omitempty"`
	Description      string `json:"description,omitempty"`
	DueDate          string `json:"due_date,omitempty"`
	SendNotification bool   `json:"send_notification,omitempty"`
}

// Create creates a new payment request (invoice).
func (s *InvoiceService) Create(ctx context.Context, req CreateInvoiceRequest, opts ...RequestOption) (*Invoice, *Response, error) {
	if req.Customer == "" {
		return nil, nil, ErrMissingCustomer
	}

	var out Invoice
	resp, err := s.client.do(ctx, "POST", "/paymentrequest", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a payment request by its numeric ID or request code.
func (s *InvoiceService) Fetch(ctx context.Context, idOrCode string, opts ...RequestOption) (*Invoice, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingReference
	}

	var out Invoice
	resp, err := s.client.do(ctx, "GET", "/paymentrequest/"+url.PathEscape(idOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// UpdateInvoiceRequest is the payload for InvoiceService.Update.
type UpdateInvoiceRequest struct {
	Amount      int64  `json:"amount,omitempty"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
}

// Update modifies a draft payment request.
func (s *InvoiceService) Update(ctx context.Context, idOrCode string, req UpdateInvoiceRequest, opts ...RequestOption) (*Invoice, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingReference
	}

	var out Invoice
	resp, err := s.client.do(ctx, "PUT", "/paymentrequest/"+url.PathEscape(idOrCode), req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Notify resends the payment request notification to the customer.
func (s *InvoiceService) Notify(ctx context.Context, idOrCode string, opts ...RequestOption) (*Response, error) {
	if idOrCode == "" {
		return nil, ErrMissingReference
	}
	return s.client.do(ctx, "POST", "/paymentrequest/notify/"+url.PathEscape(idOrCode), nil, nil, opts...)
}

// FinalizeInvoice finalizes a draft payment request so it can be sent to
// the customer.
func (s *InvoiceService) FinalizeInvoice(ctx context.Context, idOrCode string, opts ...RequestOption) (*Invoice, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingReference
	}

	var out Invoice
	resp, err := s.client.do(ctx, "POST", "/paymentrequest/finalize/"+url.PathEscape(idOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Archive archives a payment request.
func (s *InvoiceService) Archive(ctx context.Context, idOrCode string, opts ...RequestOption) (*Response, error) {
	if idOrCode == "" {
		return nil, ErrMissingReference
	}
	return s.client.do(ctx, "POST", "/paymentrequest/archive/"+url.PathEscape(idOrCode), nil, nil, opts...)
}

// ListInvoicesParams filters InvoiceService.List and InvoiceService.ListAll.
type ListInvoicesParams struct {
	Page     int
	PerPage  int
	Customer int64
	Status   string
}

func (p ListInvoicesParams) query() url.Values {
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
	return q
}

// List retrieves a single page of payment requests.
func (s *InvoiceService) List(ctx context.Context, params ListInvoicesParams, opts ...RequestOption) ([]Invoice, *Response, error) {
	path := "/paymentrequest"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Invoice
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every payment request matching params,
// fetching additional pages on demand as the caller advances it.
func (s *InvoiceService) ListAll(ctx context.Context, params ListInvoicesParams, opts ...RequestOption) *Iterator[Invoice] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Invoice, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
