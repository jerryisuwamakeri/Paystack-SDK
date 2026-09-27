package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// BulkChargeService groups the bulk-charge API operations, used to charge
// many saved authorizations in a single batch. Access it via
// Client.BulkCharges.
type BulkChargeService struct {
	client *Client
}

// BulkChargeBatch represents a Paystack bulk charge batch.
type BulkChargeBatch struct {
	ID             int64  `json:"id"`
	BatchCode      string `json:"batch_code"`
	Status         string `json:"status"`
	TotalCharges   int64  `json:"total_charges"`
	PendingCharges int64  `json:"pending_charges"`
}

// BulkChargeItem is a single charge within a bulk charge batch.
type BulkChargeItem struct {
	// Authorization is the authorization code to charge. Required.
	Authorization string `json:"authorization"`

	// Amount is the amount to charge, in the smallest currency unit.
	// Required.
	Amount int64 `json:"amount"`

	Reference string `json:"reference,omitempty"`
}

// Initiate submits a batch of charges to be processed asynchronously.
func (s *BulkChargeService) Initiate(ctx context.Context, charges []BulkChargeItem, opts ...RequestOption) (*BulkChargeBatch, *Response, error) {
	if len(charges) == 0 {
		return nil, nil, ErrMissingBulkCharges
	}

	var out BulkChargeBatch
	resp, err := s.client.do(ctx, "POST", "/bulkcharge", charges, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a bulk charge batch by its numeric ID or batch code.
func (s *BulkChargeService) Fetch(ctx context.Context, idOrCode string, opts ...RequestOption) (*BulkChargeBatch, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingBatchCode
	}

	var out BulkChargeBatch
	resp, err := s.client.do(ctx, "GET", "/bulkcharge/"+url.PathEscape(idOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// PauseBatch pauses a running bulk charge batch.
func (s *BulkChargeService) PauseBatch(ctx context.Context, batchCode string, opts ...RequestOption) (*Response, error) {
	if batchCode == "" {
		return nil, ErrMissingBatchCode
	}
	return s.client.do(ctx, "GET", "/bulkcharge/pause/"+url.PathEscape(batchCode), nil, nil, opts...)
}

// ResumeBatch resumes a paused bulk charge batch.
func (s *BulkChargeService) ResumeBatch(ctx context.Context, batchCode string, opts ...RequestOption) (*Response, error) {
	if batchCode == "" {
		return nil, ErrMissingBatchCode
	}
	return s.client.do(ctx, "GET", "/bulkcharge/resume/"+url.PathEscape(batchCode), nil, nil, opts...)
}

// ListBulkChargeBatchesParams filters BulkChargeService.List and
// BulkChargeService.ListAll.
type ListBulkChargeBatchesParams struct {
	Page    int
	PerPage int
}

func (p ListBulkChargeBatchesParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// List retrieves a single page of bulk charge batches.
func (s *BulkChargeService) List(ctx context.Context, params ListBulkChargeBatchesParams, opts ...RequestOption) ([]BulkChargeBatch, *Response, error) {
	path := "/bulkcharge"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []BulkChargeBatch
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every bulk charge batch, fetching
// additional pages on demand as the caller advances it.
func (s *BulkChargeService) ListAll(ctx context.Context, params ListBulkChargeBatchesParams, opts ...RequestOption) *Iterator[BulkChargeBatch] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]BulkChargeBatch, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
