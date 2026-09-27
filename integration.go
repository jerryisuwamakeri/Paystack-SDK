package paystack

import "context"

// IntegrationService groups integration-level configuration API
// operations. Access it via Client.Integration.
type IntegrationService struct {
	client *Client
}

// PaymentSessionTimeout is the number of seconds before an incomplete
// checkout session expires.
type PaymentSessionTimeout struct {
	PaymentSessionTimeout int `json:"payment_session_timeout"`
}

// FetchPaymentSessionTimeout retrieves the current checkout session
// timeout for the integration.
func (s *IntegrationService) FetchPaymentSessionTimeout(ctx context.Context, opts ...RequestOption) (*PaymentSessionTimeout, *Response, error) {
	var out PaymentSessionTimeout
	resp, err := s.client.do(ctx, "GET", "/integration/payment_session_timeout", nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// UpdatePaymentSessionTimeout sets the checkout session timeout, in
// seconds, for the integration.
func (s *IntegrationService) UpdatePaymentSessionTimeout(ctx context.Context, seconds int, opts ...RequestOption) (*PaymentSessionTimeout, *Response, error) {
	req := struct {
		Timeout int `json:"timeout"`
	}{Timeout: seconds}

	var out PaymentSessionTimeout
	resp, err := s.client.do(ctx, "PUT", "/integration/payment_session_timeout", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}
