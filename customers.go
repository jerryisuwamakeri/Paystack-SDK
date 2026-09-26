package paystack

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// CustomerService groups the customer-related Paystack API operations.
// Access it via Client.Customers.
type CustomerService struct {
	client *Client
}

// Customer represents a Paystack customer record.
type Customer struct {
	ID           int64          `json:"id"`
	CustomerCode string         `json:"customer_code"`
	Email        string         `json:"email"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Phone        string         `json:"phone"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
}

// CreateCustomerRequest is the payload for CustomerService.Create.
type CreateCustomerRequest struct {
	Email     string         `json:"email"`
	FirstName string         `json:"first_name,omitempty"`
	LastName  string         `json:"last_name,omitempty"`
	Phone     string         `json:"phone,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// Create registers a new customer.
func (s *CustomerService) Create(ctx context.Context, req CreateCustomerRequest, opts ...RequestOption) (*Customer, *Response, error) {
	if req.Email == "" {
		return nil, nil, ErrMissingEmail
	}

	var out Customer
	resp, err := s.client.do(ctx, "POST", "/customer", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a customer by email address or customer code.
func (s *CustomerService) Fetch(ctx context.Context, emailOrCode string, opts ...RequestOption) (*Customer, *Response, error) {
	if emailOrCode == "" {
		return nil, nil, ErrMissingReference
	}

	var out Customer
	resp, err := s.client.do(ctx, "GET", "/customer/"+url.PathEscape(emailOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// UpdateCustomerRequest is the payload for CustomerService.Update.
type UpdateCustomerRequest struct {
	FirstName string         `json:"first_name,omitempty"`
	LastName  string         `json:"last_name,omitempty"`
	Phone     string         `json:"phone,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// Update modifies an existing customer identified by customer code.
func (s *CustomerService) Update(ctx context.Context, customerCode string, req UpdateCustomerRequest, opts ...RequestOption) (*Customer, *Response, error) {
	if customerCode == "" {
		return nil, nil, ErrMissingReference
	}

	var out Customer
	resp, err := s.client.do(ctx, "PUT", "/customer/"+url.PathEscape(customerCode), req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ListCustomersParams filters CustomerService.List and
// CustomerService.ListAll.
type ListCustomersParams struct {
	Page    int
	PerPage int
}

func (p ListCustomersParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// List retrieves a single page of customers.
func (s *CustomerService) List(ctx context.Context, params ListCustomersParams, opts ...RequestOption) ([]Customer, *Response, error) {
	path := "/customer"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Customer
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every customer, fetching additional
// pages on demand as the caller advances it.
func (s *CustomerService) ListAll(ctx context.Context, params ListCustomersParams, opts ...RequestOption) *Iterator[Customer] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Customer, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
