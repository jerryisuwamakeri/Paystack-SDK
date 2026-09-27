package paystack

import (
	"context"
	"net/url"
)

// VerificationService groups identity and account verification API
// operations. Access it via Client.Verification.
type VerificationService struct {
	client *Client
}

// ResolvedAccount is returned by VerificationService.ResolveAccount.
type ResolvedAccount struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankID        int64  `json:"bank_id"`
}

// ResolveAccount looks up the account name for a given account number at a
// bank, useful for confirming transfer recipient details with the customer
// before saving them.
func (s *VerificationService) ResolveAccount(ctx context.Context, accountNumber, bankCode string, opts ...RequestOption) (*ResolvedAccount, *Response, error) {
	if accountNumber == "" {
		return nil, nil, ErrMissingAccountNumber
	}
	if bankCode == "" {
		return nil, nil, ErrMissingBankCode
	}

	q := url.Values{}
	q.Set("account_number", accountNumber)
	q.Set("bank_code", bankCode)

	var out ResolvedAccount
	resp, err := s.client.do(ctx, "GET", "/bank/resolve?"+q.Encode(), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// BVNDetails is returned by VerificationService.ResolveBVN.
type BVNDetails struct {
	BVN       string `json:"bvn"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Mobile    string `json:"mobile"`
}

// ResolveBVN retrieves the identity details associated with a Bank
// Verification Number.
func (s *VerificationService) ResolveBVN(ctx context.Context, bvn string, opts ...RequestOption) (*BVNDetails, *Response, error) {
	if bvn == "" {
		return nil, nil, ErrMissingBVN
	}

	var out BVNDetails
	resp, err := s.client.do(ctx, "GET", "/bank/resolve_bvn/"+url.PathEscape(bvn), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// CardBIN is returned by VerificationService.ResolveCardBIN.
type CardBIN struct {
	Bin         string `json:"bin"`
	Brand       string `json:"brand"`
	SubBrand    string `json:"sub_brand,omitempty"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	CardType    string `json:"card_type"`
	Bank        string `json:"bank"`
}

// ResolveCardBIN looks up the issuing bank and card type for a card's bank
// identification number, useful for routing or risk decisions before a
// charge is attempted.
func (s *VerificationService) ResolveCardBIN(ctx context.Context, bin string, opts ...RequestOption) (*CardBIN, *Response, error) {
	if bin == "" {
		return nil, nil, ErrMissingBIN
	}

	var out CardBIN
	resp, err := s.client.do(ctx, "GET", "/decision/bin/"+url.PathEscape(bin), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}
