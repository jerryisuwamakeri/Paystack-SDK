package paystack

import (
	"context"
	"net/url"
)

// ChargeService groups the direct charge API operations, used to charge a
// card, bank account, mobile money, or a saved authorization without
// redirecting the customer to a hosted checkout page. Access it via
// Client.Charge.
//
// A direct charge can pause partway through and require additional
// customer input (a card PIN, an OTP, a phone number, or a birthday)
// before completing; use the corresponding Submit* method to continue it.
type ChargeService struct {
	client *Client
}

// ChargeRequest is the payload for ChargeService.Create.
type ChargeRequest struct {
	Email             string         `json:"email"`
	Amount            int64          `json:"amount"`
	AuthorizationCode string         `json:"authorization_code,omitempty"`
	Reference         string         `json:"reference,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
}

// Create initiates a direct charge.
func (s *ChargeService) Create(ctx context.Context, req ChargeRequest, opts ...RequestOption) (*Transaction, *Response, error) {
	if req.Email == "" {
		return nil, nil, ErrMissingEmail
	}
	if req.Amount <= 0 {
		return nil, nil, ErrInvalidAmount
	}

	var out Transaction
	resp, err := s.client.do(ctx, "POST", "/charge", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// SubmitPINRequest is the payload for ChargeService.SubmitPIN.
type SubmitPINRequest struct {
	PIN       string `json:"pin"`
	Reference string `json:"reference"`
}

// SubmitPIN continues a charge that is awaiting the customer's card PIN.
func (s *ChargeService) SubmitPIN(ctx context.Context, req SubmitPINRequest, opts ...RequestOption) (*Transaction, *Response, error) {
	if req.PIN == "" {
		return nil, nil, ErrMissingPIN
	}
	if req.Reference == "" {
		return nil, nil, ErrMissingReference
	}

	var out Transaction
	resp, err := s.client.do(ctx, "POST", "/charge/submit_pin", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// SubmitOTPRequest is the payload for ChargeService.SubmitOTP.
type SubmitOTPRequest struct {
	OTP       string `json:"otp"`
	Reference string `json:"reference"`
}

// SubmitOTP continues a charge that is awaiting a one-time password.
func (s *ChargeService) SubmitOTP(ctx context.Context, req SubmitOTPRequest, opts ...RequestOption) (*Transaction, *Response, error) {
	if req.OTP == "" {
		return nil, nil, ErrMissingOTP
	}
	if req.Reference == "" {
		return nil, nil, ErrMissingReference
	}

	var out Transaction
	resp, err := s.client.do(ctx, "POST", "/charge/submit_otp", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// SubmitPhoneRequest is the payload for ChargeService.SubmitPhone.
type SubmitPhoneRequest struct {
	Phone     string `json:"phone"`
	Reference string `json:"reference"`
}

// SubmitPhone continues a charge that is awaiting the customer's phone
// number.
func (s *ChargeService) SubmitPhone(ctx context.Context, req SubmitPhoneRequest, opts ...RequestOption) (*Transaction, *Response, error) {
	if req.Phone == "" {
		return nil, nil, ErrMissingPhone
	}
	if req.Reference == "" {
		return nil, nil, ErrMissingReference
	}

	var out Transaction
	resp, err := s.client.do(ctx, "POST", "/charge/submit_phone", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// SubmitBirthdayRequest is the payload for ChargeService.SubmitBirthday.
type SubmitBirthdayRequest struct {
	Birthday  string `json:"birthday"`
	Reference string `json:"reference"`
}

// SubmitBirthday continues a charge that is awaiting the customer's date
// of birth.
func (s *ChargeService) SubmitBirthday(ctx context.Context, req SubmitBirthdayRequest, opts ...RequestOption) (*Transaction, *Response, error) {
	if req.Birthday == "" {
		return nil, nil, ErrMissingBirthday
	}
	if req.Reference == "" {
		return nil, nil, ErrMissingReference
	}

	var out Transaction
	resp, err := s.client.do(ctx, "POST", "/charge/submit_birthday", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// CheckPending polls the status of a charge left in a pending state.
func (s *ChargeService) CheckPending(ctx context.Context, reference string, opts ...RequestOption) (*Transaction, *Response, error) {
	if reference == "" {
		return nil, nil, ErrMissingReference
	}

	var out Transaction
	resp, err := s.client.do(ctx, "GET", "/charge/"+url.PathEscape(reference), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}
