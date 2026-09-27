package paystack

import "context"

// ApplePayDomainService groups the Apple Pay domain registration API
// operations. Access it via Client.ApplePayDomains.
type ApplePayDomainService struct {
	client *Client
}

type applePayDomainRequest struct {
	DomainName string `json:"domainName"`
}

// Register registers a domain for use with Apple Pay.
func (s *ApplePayDomainService) Register(ctx context.Context, domainName string, opts ...RequestOption) (*Response, error) {
	if domainName == "" {
		return nil, ErrMissingDomainName
	}
	return s.client.do(ctx, "POST", "/apple-pay/domain", applePayDomainRequest{DomainName: domainName}, nil, opts...)
}

// List retrieves the domains currently registered for Apple Pay.
func (s *ApplePayDomainService) List(ctx context.Context, opts ...RequestOption) ([]string, *Response, error) {
	var out struct {
		DomainNames []string `json:"domainNames"`
	}
	resp, err := s.client.do(ctx, "GET", "/apple-pay/domain", nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out.DomainNames, resp, nil
}

// Unregister removes a domain from Apple Pay.
func (s *ApplePayDomainService) Unregister(ctx context.Context, domainName string, opts ...RequestOption) (*Response, error) {
	if domainName == "" {
		return nil, ErrMissingDomainName
	}
	return s.client.do(ctx, "DELETE", "/apple-pay/domain", applePayDomainRequest{DomainName: domainName}, nil, opts...)
}
