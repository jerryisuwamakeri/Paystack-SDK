package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// PaymentPageService groups the payment-page API operations, used for
// no-code hosted checkout links. Access it via Client.PaymentPages.
type PaymentPageService struct {
	client *Client
}

// PaymentPage represents a Paystack hosted payment page.
type PaymentPage struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Slug        string `json:"slug"`
	Amount      int64  `json:"amount,omitempty"`
	Currency    string `json:"currency,omitempty"`
	Active      bool   `json:"active"`
}

// CreatePaymentPageRequest is the payload for PaymentPageService.Create.
type CreatePaymentPageRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// Amount is optional; omit it to let the customer choose how much to
	// pay.
	Amount int64  `json:"amount,omitempty"`
	Slug   string `json:"slug,omitempty"`
}

// Create creates a new hosted payment page.
func (s *PaymentPageService) Create(ctx context.Context, req CreatePaymentPageRequest, opts ...RequestOption) (*PaymentPage, *Response, error) {
	if req.Name == "" {
		return nil, nil, ErrMissingName
	}

	var out PaymentPage
	resp, err := s.client.do(ctx, "POST", "/page", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a payment page by its numeric ID or slug.
func (s *PaymentPageService) Fetch(ctx context.Context, idOrSlug string, opts ...RequestOption) (*PaymentPage, *Response, error) {
	if idOrSlug == "" {
		return nil, nil, ErrMissingReference
	}

	var out PaymentPage
	resp, err := s.client.do(ctx, "GET", "/page/"+url.PathEscape(idOrSlug), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// UpdatePaymentPageRequest is the payload for PaymentPageService.Update.
type UpdatePaymentPageRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Amount      int64  `json:"amount,omitempty"`
	Active      *bool  `json:"active,omitempty"`
}

// Update modifies an existing payment page.
func (s *PaymentPageService) Update(ctx context.Context, idOrSlug string, req UpdatePaymentPageRequest, opts ...RequestOption) (*PaymentPage, *Response, error) {
	if idOrSlug == "" {
		return nil, nil, ErrMissingReference
	}

	var out PaymentPage
	resp, err := s.client.do(ctx, "PUT", "/page/"+url.PathEscape(idOrSlug), req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ListPaymentPagesParams filters PaymentPageService.List and
// PaymentPageService.ListAll.
type ListPaymentPagesParams struct {
	Page    int
	PerPage int
}

func (p ListPaymentPagesParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// List retrieves a single page of payment pages.
func (s *PaymentPageService) List(ctx context.Context, params ListPaymentPagesParams, opts ...RequestOption) ([]PaymentPage, *Response, error) {
	path := "/page"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []PaymentPage
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every payment page, fetching additional
// pages on demand as the caller advances it.
func (s *PaymentPageService) ListAll(ctx context.Context, params ListPaymentPagesParams, opts ...RequestOption) *Iterator[PaymentPage] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]PaymentPage, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
