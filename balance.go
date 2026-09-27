package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// BalanceService groups the balance and balance-ledger API operations.
// Access it via Client.Balance.
type BalanceService struct {
	client *Client
}

// Balance is the available balance for a single currency in the
// integration's account.
type Balance struct {
	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`
}

// Check retrieves the current balance for each currency in the account.
func (s *BalanceService) Check(ctx context.Context, opts ...RequestOption) ([]Balance, *Response, error) {
	var out []Balance
	resp, err := s.client.do(ctx, "GET", "/balance", nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// BalanceLedgerEntry is a single entry in the balance ledger, recording a
// change in the account balance and why it happened.
type BalanceLedgerEntry struct {
	Integration      int64  `json:"integration"`
	Domain           string `json:"domain"`
	Balance          int64  `json:"balance"`
	Currency         string `json:"currency"`
	Difference       int64  `json:"difference"`
	Reason           string `json:"reason"`
	ModelResponsible string `json:"model_responsible,omitempty"`
}

// ListBalanceLedgerParams filters BalanceService.FetchLedger and
// BalanceService.FetchLedgerAll.
type ListBalanceLedgerParams struct {
	Page    int
	PerPage int
}

func (p ListBalanceLedgerParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// FetchLedger retrieves a single page of balance ledger entries.
func (s *BalanceService) FetchLedger(ctx context.Context, params ListBalanceLedgerParams, opts ...RequestOption) ([]BalanceLedgerEntry, *Response, error) {
	path := "/balance/ledger"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []BalanceLedgerEntry
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// FetchLedgerAll returns an Iterator over every balance ledger entry,
// fetching additional pages on demand as the caller advances it.
func (s *BalanceService) FetchLedgerAll(ctx context.Context, params ListBalanceLedgerParams, opts ...RequestOption) *Iterator[BalanceLedgerEntry] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]BalanceLedgerEntry, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.FetchLedger(ctx, p, opts...)
	})
}
