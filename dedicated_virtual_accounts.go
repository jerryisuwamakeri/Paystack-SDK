package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// DedicatedVirtualAccountService groups the dedicated virtual account API
// operations. Access it via Client.DedicatedVirtualAccounts.
type DedicatedVirtualAccountService struct {
	client *Client
}

// DedicatedVirtualAccountBank identifies the bank backing a
// DedicatedVirtualAccount.
type DedicatedVirtualAccountBank struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// DedicatedVirtualAccount represents a customer-specific virtual bank
// account used to accept direct bank transfers.
type DedicatedVirtualAccount struct {
	ID            int64                       `json:"id"`
	AccountName   string                      `json:"account_name"`
	AccountNumber string                      `json:"account_number"`
	Active        bool                        `json:"active"`
	Bank          DedicatedVirtualAccountBank `json:"bank"`
}

// CreateDedicatedVirtualAccountRequest is the payload for
// DedicatedVirtualAccountService.Create.
type CreateDedicatedVirtualAccountRequest struct {
	// Customer is the customer code to assign the account to. Required.
	Customer string `json:"customer"`

	// PreferredBank restricts assignment to a specific bank, when the
	// integration supports multiple providers. Optional.
	PreferredBank string `json:"preferred_bank,omitempty"`
}

// Create assigns a new dedicated virtual account to a customer.
func (s *DedicatedVirtualAccountService) Create(ctx context.Context, req CreateDedicatedVirtualAccountRequest, opts ...RequestOption) (*DedicatedVirtualAccount, *Response, error) {
	if req.Customer == "" {
		return nil, nil, ErrMissingCustomer
	}

	var out DedicatedVirtualAccount
	resp, err := s.client.do(ctx, "POST", "/dedicated_account", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a dedicated virtual account by its numeric ID.
func (s *DedicatedVirtualAccountService) Fetch(ctx context.Context, id string, opts ...RequestOption) (*DedicatedVirtualAccount, *Response, error) {
	if id == "" {
		return nil, nil, ErrMissingReference
	}

	var out DedicatedVirtualAccount
	resp, err := s.client.do(ctx, "GET", "/dedicated_account/"+url.PathEscape(id), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Deactivate deactivates a dedicated virtual account by its numeric ID.
func (s *DedicatedVirtualAccountService) Deactivate(ctx context.Context, id string, opts ...RequestOption) (*Response, error) {
	if id == "" {
		return nil, ErrMissingReference
	}
	return s.client.do(ctx, "DELETE", "/dedicated_account/"+url.PathEscape(id), nil, nil, opts...)
}

// ListDedicatedVirtualAccountsParams filters
// DedicatedVirtualAccountService.List and
// DedicatedVirtualAccountService.ListAll.
type ListDedicatedVirtualAccountsParams struct {
	Page     int
	PerPage  int
	Active   *bool
	Currency string
}

func (p ListDedicatedVirtualAccountsParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	if p.Active != nil {
		q.Set("active", strconv.FormatBool(*p.Active))
	}
	if p.Currency != "" {
		q.Set("currency", p.Currency)
	}
	return q
}

// List retrieves a single page of dedicated virtual accounts.
func (s *DedicatedVirtualAccountService) List(ctx context.Context, params ListDedicatedVirtualAccountsParams, opts ...RequestOption) ([]DedicatedVirtualAccount, *Response, error) {
	path := "/dedicated_account"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []DedicatedVirtualAccount
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every dedicated virtual account matching
// params, fetching additional pages on demand as the caller advances it.
func (s *DedicatedVirtualAccountService) ListAll(ctx context.Context, params ListDedicatedVirtualAccountsParams, opts ...RequestOption) *Iterator[DedicatedVirtualAccount] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]DedicatedVirtualAccount, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}

// DedicatedVirtualAccountProvider is a bank capable of issuing dedicated
// virtual accounts.
type DedicatedVirtualAccountProvider struct {
	ProviderSlug string `json:"provider_slug"`
	BankID       int64  `json:"bank_id"`
	BankName     string `json:"bank_name"`
}

// ListProviders retrieves the banks available for issuing dedicated
// virtual accounts.
func (s *DedicatedVirtualAccountService) ListProviders(ctx context.Context, opts ...RequestOption) ([]DedicatedVirtualAccountProvider, *Response, error) {
	var out []DedicatedVirtualAccountProvider
	resp, err := s.client.do(ctx, "GET", "/dedicated_account/available_providers", nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}
