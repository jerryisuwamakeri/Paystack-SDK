package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// TransferRecipientService groups transfer-recipient API operations.
// Access it via Client.TransferRecipients.
type TransferRecipientService struct {
	client *Client
}

// TransferRecipientDetails carries the destination account details of a
// TransferRecipient.
type TransferRecipientDetails struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankCode      string `json:"bank_code"`
	BankName      string `json:"bank_name"`
}

// TransferRecipient represents a saved transfer destination.
type TransferRecipient struct {
	ID            int64                    `json:"id"`
	RecipientCode string                   `json:"recipient_code"`
	Type          string                   `json:"type"`
	Name          string                   `json:"name"`
	Currency      string                   `json:"currency"`
	Active        bool                     `json:"active"`
	Details       TransferRecipientDetails `json:"details"`
}

// CreateTransferRecipientRequest is the payload for
// TransferRecipientService.Create.
type CreateTransferRecipientRequest struct {
	// Type is the recipient type, for example "nuban" for Nigerian bank
	// accounts. Defaults to "nuban" when empty.
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	Currency      string `json:"currency,omitempty"`
}

// Create saves a new transfer recipient that can later be used as the
// destination of a transfer.
func (s *TransferRecipientService) Create(ctx context.Context, req CreateTransferRecipientRequest, opts ...RequestOption) (*TransferRecipient, *Response, error) {
	if req.Name == "" {
		return nil, nil, ErrMissingName
	}
	if req.AccountNumber == "" {
		return nil, nil, ErrMissingAccountNumber
	}
	if req.BankCode == "" {
		return nil, nil, ErrMissingBankCode
	}
	if req.Type == "" {
		req.Type = "nuban"
	}

	var out TransferRecipient
	resp, err := s.client.do(ctx, "POST", "/transferrecipient", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a transfer recipient by its numeric ID or recipient code.
func (s *TransferRecipientService) Fetch(ctx context.Context, idOrCode string, opts ...RequestOption) (*TransferRecipient, *Response, error) {
	if idOrCode == "" {
		return nil, nil, ErrMissingRecipientCode
	}

	var out TransferRecipient
	resp, err := s.client.do(ctx, "GET", "/transferrecipient/"+url.PathEscape(idOrCode), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ListTransferRecipientsParams filters TransferRecipientService.List and
// TransferRecipientService.ListAll.
type ListTransferRecipientsParams struct {
	Page    int
	PerPage int
}

func (p ListTransferRecipientsParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// List retrieves a single page of transfer recipients.
func (s *TransferRecipientService) List(ctx context.Context, params ListTransferRecipientsParams, opts ...RequestOption) ([]TransferRecipient, *Response, error) {
	path := "/transferrecipient"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []TransferRecipient
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every transfer recipient, fetching
// additional pages on demand as the caller advances it.
func (s *TransferRecipientService) ListAll(ctx context.Context, params ListTransferRecipientsParams, opts ...RequestOption) *Iterator[TransferRecipient] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]TransferRecipient, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
