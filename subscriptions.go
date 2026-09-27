package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// SubscriptionService groups the subscription-related Paystack API
// operations. Access it via Client.Subscriptions.
type SubscriptionService struct {
	client *Client
}

// Subscription represents a Paystack subscription record.
type Subscription struct {
	ID               int64  `json:"id"`
	SubscriptionCode string `json:"subscription_code"`
	EmailToken       string `json:"email_token"`
	Amount           int64  `json:"amount"`
	Status           string `json:"status"`
	Plan             *Plan  `json:"plan,omitempty"`
}

// CreateSubscriptionRequest is the payload for SubscriptionService.Create.
type CreateSubscriptionRequest struct {
	// Customer is the customer's email address or customer code. Required.
	Customer string `json:"customer"`

	// Plan is the plan code to subscribe the customer to. Required.
	Plan string `json:"plan"`

	// Authorization is the authorization code to charge for this
	// subscription. Optional; defaults to the customer's most recent
	// authorization when omitted.
	Authorization string `json:"authorization,omitempty"`
}

// Create subscribes a customer to a plan.
func (s *SubscriptionService) Create(ctx context.Context, req CreateSubscriptionRequest, opts ...RequestOption) (*Subscription, *Response, error) {
	if req.Customer == "" {
		return nil, nil, ErrMissingCustomer
	}
	if req.Plan == "" {
		return nil, nil, ErrMissingPlanCode
	}

	var out Subscription
	resp, err := s.client.do(ctx, "POST", "/subscription", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a subscription by its numeric ID or subscription code.
func (s *SubscriptionService) Fetch(ctx context.Context, idOrCode string, opts ...RequestOption) (*Subscription, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingSubscriptionCode
	}

	var out Subscription
	resp, err := s.client.do(ctx, "GET", "/subscription/"+url.PathEscape(idOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ManageSubscriptionRequest is the payload for SubscriptionService.Enable
// and SubscriptionService.Disable.
type ManageSubscriptionRequest struct {
	Code  string `json:"code"`
	Token string `json:"token"`
}

func (req ManageSubscriptionRequest) validate() error {
	if req.Code == "" {
		return ErrMissingSubscriptionCode
	}
	if req.Token == "" {
		return ErrMissingEmailToken
	}
	return nil
}

// Enable reactivates a subscription that is not currently active.
func (s *SubscriptionService) Enable(ctx context.Context, req ManageSubscriptionRequest, opts ...RequestOption) (*Response, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return s.client.do(ctx, "POST", "/subscription/enable", req, nil, opts...)
}

// Disable cancels an active subscription.
func (s *SubscriptionService) Disable(ctx context.Context, req ManageSubscriptionRequest, opts ...RequestOption) (*Response, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return s.client.do(ctx, "POST", "/subscription/disable", req, nil, opts...)
}

// UpdateSubscriptionLink is returned by
// SubscriptionService.GenerateUpdateLink.
type UpdateSubscriptionLink struct {
	Link string `json:"link"`
}

// GenerateUpdateLink creates a link the customer can use to update the
// card on file for a subscription.
func (s *SubscriptionService) GenerateUpdateLink(ctx context.Context, code string, opts ...RequestOption) (*UpdateSubscriptionLink, *Response, error) {
	if code == "" {
		return nil, nil, ErrMissingSubscriptionCode
	}

	var out UpdateSubscriptionLink
	resp, err := s.client.do(ctx, "GET", "/subscription/"+url.PathEscape(code)+"/manage/link", nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// SendUpdateLink emails the customer a link to update the card on file for
// a subscription.
func (s *SubscriptionService) SendUpdateLink(ctx context.Context, code string, opts ...RequestOption) (*Response, error) {
	if code == "" {
		return nil, ErrMissingSubscriptionCode
	}
	return s.client.do(ctx, "POST", "/subscription/"+url.PathEscape(code)+"/manage/email", nil, nil, opts...)
}

// ListSubscriptionsParams filters SubscriptionService.List and
// SubscriptionService.ListAll.
type ListSubscriptionsParams struct {
	Page     int
	PerPage  int
	Customer int64
	Plan     string
}

func (p ListSubscriptionsParams) query() url.Values {
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
	if p.Plan != "" {
		q.Set("plan", p.Plan)
	}
	return q
}

// List retrieves a single page of subscriptions.
func (s *SubscriptionService) List(ctx context.Context, params ListSubscriptionsParams, opts ...RequestOption) ([]Subscription, *Response, error) {
	path := "/subscription"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Subscription
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every subscription matching params,
// fetching additional pages on demand as the caller advances it.
func (s *SubscriptionService) ListAll(ctx context.Context, params ListSubscriptionsParams, opts ...RequestOption) *Iterator[Subscription] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Subscription, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
