package paystack

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// SettlementService groups the settlement-related Paystack API operations.
// Access it via Client.Settlements.
type SettlementService struct {
	client *Client
}

// Settlement represents a Paystack settlement payout record.
type Settlement struct {
	ID             int64     `json:"id"`
	Status         string    `json:"status"`
	Currency       string    `json:"currency"`
	TotalAmount    int64     `json:"total_amount"`
	TotalFees      int64     `json:"total_fees"`
	SettlementDate time.Time `json:"settlement_date"`
}

// ListSettlementsParams filters SettlementService.List and
// SettlementService.ListAll.
type ListSettlementsParams struct {
	Page    int
	PerPage int
	Status  string
	From    time.Time
	To      time.Time
}

func (p ListSettlementsParams) query() url.Values {
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
	if !p.From.IsZero() {
		q.Set("from", p.From.Format(time.RFC3339))
	}
	if !p.To.IsZero() {
		q.Set("to", p.To.Format(time.RFC3339))
	}
	return q
}

// List retrieves a single page of settlements.
func (s *SettlementService) List(ctx context.Context, params ListSettlementsParams, opts ...RequestOption) ([]Settlement, *Response, error) {
	path := "/settlement"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Settlement
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every settlement matching params,
// fetching additional pages on demand as the caller advances it.
func (s *SettlementService) ListAll(ctx context.Context, params ListSettlementsParams, opts ...RequestOption) *Iterator[Settlement] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Settlement, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
