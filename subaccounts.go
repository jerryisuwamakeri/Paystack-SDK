package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// SubaccountService groups the subaccount-related Paystack API operations,
// used for splitting settlement payouts with other parties. Access it via
// Client.Subaccounts.
type SubaccountService struct {
	client *Client
}

// Subaccount represents a Paystack subaccount.
type Subaccount struct {
	ID               int64   `json:"id"`
	SubaccountCode   string  `json:"subaccount_code"`
	BusinessName     string  `json:"business_name"`
	SettlementBank   string  `json:"settlement_bank"`
	AccountNumber    string  `json:"account_number"`
	PercentageCharge float64 `json:"percentage_charge"`
	Active           bool    `json:"active"`
}

// CreateSubaccountRequest is the payload for SubaccountService.Create.
type CreateSubaccountRequest struct {
	BusinessName     string  `json:"business_name"`
	SettlementBank   string  `json:"settlement_bank"`
	AccountNumber    string  `json:"account_number"`
	PercentageCharge float64 `json:"percentage_charge"`
}

// Create registers a new subaccount.
func (s *SubaccountService) Create(ctx context.Context, req CreateSubaccountRequest, opts ...RequestOption) (*Subaccount, *Response, error) {
	if req.BusinessName == "" {
		return nil, nil, ErrMissingBusinessName
	}
	if req.SettlementBank == "" {
		return nil, nil, ErrMissingSettlementBank
	}
	if req.AccountNumber == "" {
		return nil, nil, ErrMissingAccountNumber
	}

	var out Subaccount
	resp, err := s.client.do(ctx, "POST", "/subaccount", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a subaccount by its numeric ID or subaccount code.
func (s *SubaccountService) Fetch(ctx context.Context, idOrCode string, opts ...RequestOption) (*Subaccount, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingSubaccountCode
	}

	var out Subaccount
	resp, err := s.client.do(ctx, "GET", "/subaccount/"+url.PathEscape(idOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// UpdateSubaccountRequest is the payload for SubaccountService.Update.
type UpdateSubaccountRequest struct {
	BusinessName     string  `json:"business_name,omitempty"`
	SettlementBank   string  `json:"settlement_bank,omitempty"`
	AccountNumber    string  `json:"account_number,omitempty"`
	PercentageCharge float64 `json:"percentage_charge,omitempty"`
	Active           *bool   `json:"active,omitempty"`
}

// Update modifies an existing subaccount.
func (s *SubaccountService) Update(ctx context.Context, idOrCode string, req UpdateSubaccountRequest, opts ...RequestOption) (*Subaccount, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingSubaccountCode
	}

	var out Subaccount
	resp, err := s.client.do(ctx, "PUT", "/subaccount/"+url.PathEscape(idOrCode), req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ListSubaccountsParams filters SubaccountService.List and
// SubaccountService.ListAll.
type ListSubaccountsParams struct {
	Page    int
	PerPage int
}

func (p ListSubaccountsParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// List retrieves a single page of subaccounts.
func (s *SubaccountService) List(ctx context.Context, params ListSubaccountsParams, opts ...RequestOption) ([]Subaccount, *Response, error) {
	path := "/subaccount"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Subaccount
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every subaccount, fetching additional
// pages on demand as the caller advances it.
func (s *SubaccountService) ListAll(ctx context.Context, params ListSubaccountsParams, opts ...RequestOption) *Iterator[Subaccount] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Subaccount, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
