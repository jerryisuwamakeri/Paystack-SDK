package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// MiscellaneousService groups Paystack reference-data API operations, such
// as listing supported banks and countries. Access it via
// Client.Miscellaneous.
type MiscellaneousService struct {
	client *Client
}

// Bank represents a bank supported by Paystack.
type Bank struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Code     string `json:"code"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
}

// ListBanksParams filters MiscellaneousService.ListBanks.
type ListBanksParams struct {
	Country  string
	Currency string
	Page     int
	PerPage  int
}

func (p ListBanksParams) query() url.Values {
	q := url.Values{}
	if p.Country != "" {
		q.Set("country", p.Country)
	}
	if p.Currency != "" {
		q.Set("currency", p.Currency)
	}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// ListBanks retrieves the list of banks supported by Paystack, optionally
// filtered by country or currency.
func (s *MiscellaneousService) ListBanks(ctx context.Context, params ListBanksParams, opts ...RequestOption) ([]Bank, *Response, error) {
	path := "/bank"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Bank
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// Country represents a country supported by Paystack.
type Country struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name"`
	ISOCode             string `json:"iso_code"`
	DefaultCurrencyCode string `json:"default_currency_code"`
}

// ListCountries retrieves the list of countries supported by Paystack.
func (s *MiscellaneousService) ListCountries(ctx context.Context, opts ...RequestOption) ([]Country, *Response, error) {
	var out []Country
	resp, err := s.client.do(ctx, "GET", "/country", nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}
