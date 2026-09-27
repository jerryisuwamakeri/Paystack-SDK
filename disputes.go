package paystack

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// DisputeService groups the dispute (chargeback) API operations. Access it
// via Client.Disputes.
type DisputeService struct {
	client *Client
}

// Dispute represents a Paystack transaction dispute.
type Dispute struct {
	ID         int64  `json:"id"`
	Status     string `json:"status"`
	Category   string `json:"category,omitempty"`
	Currency   string `json:"currency"`
	Amount     int64  `json:"amount"`
	Resolution string `json:"resolution,omitempty"`
}

// ListDisputesParams filters DisputeService.List and DisputeService.ListAll.
type ListDisputesParams struct {
	Page        int
	PerPage     int
	Status      string
	Transaction string
	From        time.Time
	To          time.Time
}

func (p ListDisputesParams) query() url.Values {
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
	if p.Transaction != "" {
		q.Set("transaction", p.Transaction)
	}
	if !p.From.IsZero() {
		q.Set("from", p.From.Format(time.RFC3339))
	}
	if !p.To.IsZero() {
		q.Set("to", p.To.Format(time.RFC3339))
	}
	return q
}

// List retrieves a single page of disputes.
func (s *DisputeService) List(ctx context.Context, params ListDisputesParams, opts ...RequestOption) ([]Dispute, *Response, error) {
	path := "/dispute"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Dispute
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every dispute matching params, fetching
// additional pages on demand as the caller advances it.
func (s *DisputeService) ListAll(ctx context.Context, params ListDisputesParams, opts ...RequestOption) *Iterator[Dispute] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Dispute, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}

// Fetch retrieves a dispute by its numeric ID.
func (s *DisputeService) Fetch(ctx context.Context, id string, opts ...RequestOption) (*Dispute, *Response, error) {
	if id == "" {
		return nil, nil, ErrMissingReference
	}

	var out Dispute
	resp, err := s.client.do(ctx, "GET", "/dispute/"+url.PathEscape(id), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// AddDisputeEvidenceRequest is the payload for DisputeService.AddEvidence.
type AddDisputeEvidenceRequest struct {
	CustomerEmail   string `json:"customer_email"`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	ServiceDetails  string `json:"service_details"`
	DeliveryAddress string `json:"delivery_address,omitempty"`
	DeliveryDate    string `json:"delivery_date,omitempty"`
}

// AddEvidence submits evidence for a dispute.
func (s *DisputeService) AddEvidence(ctx context.Context, id string, req AddDisputeEvidenceRequest, opts ...RequestOption) (*Response, error) {
	if id == "" {
		return nil, ErrMissingReference
	}
	if req.CustomerEmail == "" {
		return nil, ErrMissingEmail
	}
	if req.CustomerName == "" {
		return nil, ErrMissingCustomerName
	}
	if req.CustomerPhone == "" {
		return nil, ErrMissingCustomerPhone
	}
	if req.ServiceDetails == "" {
		return nil, ErrMissingServiceDetails
	}

	return s.client.do(ctx, "POST", "/dispute/"+url.PathEscape(id)+"/evidence", req, nil, opts...)
}

// ResolveDisputeRequest is the payload for DisputeService.Resolve.
type ResolveDisputeRequest struct {
	// Resolution is either "merchant-accepted" or "declined". Required.
	Resolution string `json:"resolution"`

	Message          string `json:"message"`
	UploadedFilename string `json:"uploaded_filename"`
	RefundAmount     int64  `json:"refund_amount,omitempty"`
	Evidence         int64  `json:"evidence,omitempty"`
}

// Resolve finalizes a dispute with the merchant's decision.
func (s *DisputeService) Resolve(ctx context.Context, id string, req ResolveDisputeRequest, opts ...RequestOption) (*Dispute, *Response, error) {
	if id == "" {
		return nil, nil, ErrMissingReference
	}
	if req.Resolution == "" {
		return nil, nil, ErrMissingResolution
	}
	if req.Message == "" {
		return nil, nil, ErrMissingMessage
	}
	if req.UploadedFilename == "" {
		return nil, nil, ErrMissingUploadedFilename
	}

	var out Dispute
	resp, err := s.client.do(ctx, "PUT", "/dispute/"+url.PathEscape(id)+"/resolve", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}
