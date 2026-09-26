package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// SplitPaymentService groups the transaction-split API operations, used to
// share a single payment's proceeds across multiple subaccounts. Access it
// via Client.SplitPayments.
type SplitPaymentService struct {
	client *Client
}

// SplitShare assigns a subaccount's portion of a split payment.
type SplitShare struct {
	Subaccount string  `json:"subaccount"`
	Share      float64 `json:"share"`
}

// SplitPayment represents a Paystack transaction split configuration.
type SplitPayment struct {
	ID          int64        `json:"id"`
	SplitCode   string       `json:"split_code"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Currency    string       `json:"currency"`
	Active      bool         `json:"active"`
	Subaccounts []SplitShare `json:"subaccounts"`
}

// CreateSplitPaymentRequest is the payload for SplitPaymentService.Create.
type CreateSplitPaymentRequest struct {
	Name        string       `json:"name"`
	Type        string       `json:"type"` // "percentage" or "flat"
	Currency    string       `json:"currency"`
	Subaccounts []SplitShare `json:"subaccounts"`
	BearerType  string       `json:"bearer_type,omitempty"`
}

// Create creates a new transaction split configuration.
func (s *SplitPaymentService) Create(ctx context.Context, req CreateSplitPaymentRequest, opts ...RequestOption) (*SplitPayment, *Response, error) {
	if req.Name == "" {
		return nil, nil, ErrMissingName
	}
	if req.Type == "" {
		return nil, nil, ErrMissingSplitType
	}
	if len(req.Subaccounts) == 0 {
		return nil, nil, ErrMissingSplitSubaccounts
	}

	var out SplitPayment
	resp, err := s.client.do(ctx, "POST", "/split", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a split payment by its numeric ID or split code.
func (s *SplitPaymentService) Fetch(ctx context.Context, idOrCode string, opts ...RequestOption) (*SplitPayment, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingSplitCode
	}

	var out SplitPayment
	resp, err := s.client.do(ctx, "GET", "/split/"+url.PathEscape(idOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// UpdateSplitPaymentRequest is the payload for SplitPaymentService.Update.
type UpdateSplitPaymentRequest struct {
	Name   string `json:"name,omitempty"`
	Active *bool  `json:"active,omitempty"`
}

// Update modifies an existing split payment configuration.
func (s *SplitPaymentService) Update(ctx context.Context, idOrCode string, req UpdateSplitPaymentRequest, opts ...RequestOption) (*SplitPayment, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingSplitCode
	}

	var out SplitPayment
	resp, err := s.client.do(ctx, "PUT", "/split/"+url.PathEscape(idOrCode), req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ListSplitPaymentsParams filters SplitPaymentService.List and
// SplitPaymentService.ListAll.
type ListSplitPaymentsParams struct {
	Page    int
	PerPage int
	Name    string
	Active  *bool
}

func (p ListSplitPaymentsParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	if p.Name != "" {
		q.Set("name", p.Name)
	}
	if p.Active != nil {
		q.Set("active", strconv.FormatBool(*p.Active))
	}
	return q
}

// List retrieves a single page of split payments.
func (s *SplitPaymentService) List(ctx context.Context, params ListSplitPaymentsParams, opts ...RequestOption) ([]SplitPayment, *Response, error) {
	path := "/split"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []SplitPayment
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every split payment matching params,
// fetching additional pages on demand as the caller advances it.
func (s *SplitPaymentService) ListAll(ctx context.Context, params ListSplitPaymentsParams, opts ...RequestOption) *Iterator[SplitPayment] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]SplitPayment, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
