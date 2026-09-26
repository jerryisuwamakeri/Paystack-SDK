package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// PlanService groups the plan-related Paystack API operations. Access it
// via Client.Plans.
type PlanService struct {
	client *Client
}

// Plan represents a Paystack subscription plan.
type Plan struct {
	ID          int64  `json:"id"`
	PlanCode    string `json:"plan_code"`
	Name        string `json:"name"`
	Amount      int64  `json:"amount"`
	Interval    string `json:"interval"`
	Currency    string `json:"currency"`
	Description string `json:"description,omitempty"`
}

// CreatePlanRequest is the payload for PlanService.Create.
type CreatePlanRequest struct {
	Name        string `json:"name"`
	Amount      int64  `json:"amount"`
	Interval    string `json:"interval"`
	Currency    string `json:"currency,omitempty"`
	Description string `json:"description,omitempty"`
}

// Create creates a new subscription plan.
func (s *PlanService) Create(ctx context.Context, req CreatePlanRequest, opts ...RequestOption) (*Plan, *Response, error) {
	if req.Name == "" {
		return nil, nil, ErrMissingName
	}
	if req.Amount <= 0 {
		return nil, nil, ErrInvalidAmount
	}
	if req.Interval == "" {
		return nil, nil, ErrMissingInterval
	}

	var out Plan
	resp, err := s.client.do(ctx, "POST", "/plan", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a plan by its numeric ID or plan code.
func (s *PlanService) Fetch(ctx context.Context, idOrCode string, opts ...RequestOption) (*Plan, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingPlanCode
	}

	var out Plan
	resp, err := s.client.do(ctx, "GET", "/plan/"+url.PathEscape(idOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// UpdatePlanRequest is the payload for PlanService.Update. Zero-value
// fields are omitted from the request and left unchanged on the plan.
type UpdatePlanRequest struct {
	Name        string `json:"name,omitempty"`
	Amount      int64  `json:"amount,omitempty"`
	Interval    string `json:"interval,omitempty"`
	Description string `json:"description,omitempty"`
}

// Update modifies an existing plan.
func (s *PlanService) Update(ctx context.Context, idOrCode string, req UpdatePlanRequest, opts ...RequestOption) (*Response, error) {
	if idOrCode == "" {
		return nil, ErrMissingPlanCode
	}

	resp, err := s.client.do(ctx, "PUT", "/plan/"+url.PathEscape(idOrCode), req, nil, opts...)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListPlansParams filters PlanService.List and PlanService.ListAll.
type ListPlansParams struct {
	Page    int
	PerPage int
}

func (p ListPlansParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// List retrieves a single page of plans.
func (s *PlanService) List(ctx context.Context, params ListPlansParams, opts ...RequestOption) ([]Plan, *Response, error) {
	path := "/plan"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Plan
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every plan, fetching additional pages
// on demand as the caller advances it.
func (s *PlanService) ListAll(ctx context.Context, params ListPlansParams, opts ...RequestOption) *Iterator[Plan] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Plan, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
