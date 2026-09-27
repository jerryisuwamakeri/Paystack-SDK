package paystack

import (
	"fmt"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.paystack.co"
	defaultTimeout = 30 * time.Second
)

// Client is a Paystack API client. Once constructed by NewClient, a
// Client's configuration is immutable and it is safe for concurrent use by
// multiple goroutines: callers may share a single Client across an entire
// application.
type Client struct {
	secretKey   string
	baseURL     string
	httpClient  *http.Client
	retryPolicy RetryPolicy
	observer    TransportObserver

	// Transactions groups transaction-related API operations.
	Transactions *TransactionService

	// Customers groups customer-related API operations.
	Customers *CustomerService

	// TransferRecipients groups transfer-recipient API operations.
	TransferRecipients *TransferRecipientService

	// Transfers groups transfer-related API operations.
	Transfers *TransferService

	// Plans groups subscription plan API operations.
	Plans *PlanService

	// Subscriptions groups subscription-related API operations.
	Subscriptions *SubscriptionService

	// Refunds groups refund-related API operations.
	Refunds *RefundService

	// Miscellaneous groups reference-data API operations such as banks and
	// countries.
	Miscellaneous *MiscellaneousService

	// Verification groups identity and account verification API
	// operations.
	Verification *VerificationService

	// Products groups product-related API operations.
	Products *ProductService

	// Subaccounts groups subaccount-related API operations.
	Subaccounts *SubaccountService

	// Settlements groups settlement-related API operations.
	Settlements *SettlementService

	// DedicatedVirtualAccounts groups dedicated virtual account API
	// operations.
	DedicatedVirtualAccounts *DedicatedVirtualAccountService

	// SplitPayments groups transaction-split API operations.
	SplitPayments *SplitPaymentService

	// PaymentPages groups hosted payment page API operations.
	PaymentPages *PaymentPageService

	// BulkCharges groups bulk-charge API operations.
	BulkCharges *BulkChargeService

	// Disputes groups dispute (chargeback) API operations.
	Disputes *DisputeService

	// ApplePayDomains groups Apple Pay domain registration API operations.
	ApplePayDomains *ApplePayDomainService
}

// NewClient constructs a Client using the given secret key and options.
//
// The secret key must be supplied through runtime configuration, such as
// an environment variable:
//
//	client, err := paystack.NewClient(os.Getenv("PAYSTACK_SECRET_KEY"))
//
// The SDK never persists, logs, or otherwise exposes the secret key.
func NewClient(secretKey string, opts ...Option) (*Client, error) {
	if secretKey == "" {
		return nil, ErrMissingSecretKey
	}

	c := &Client{
		secretKey:   secretKey,
		baseURL:     defaultBaseURL,
		httpClient:  &http.Client{Timeout: defaultTimeout},
		retryPolicy: DefaultRetryPolicy(),
		observer:    noopObserver{},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Transactions = &TransactionService{client: c}
	c.Customers = &CustomerService{client: c}
	c.TransferRecipients = &TransferRecipientService{client: c}
	c.Transfers = &TransferService{client: c}
	c.Plans = &PlanService{client: c}
	c.Subscriptions = &SubscriptionService{client: c}
	c.Refunds = &RefundService{client: c}
	c.Miscellaneous = &MiscellaneousService{client: c}
	c.Verification = &VerificationService{client: c}
	c.Products = &ProductService{client: c}
	c.Subaccounts = &SubaccountService{client: c}
	c.Settlements = &SettlementService{client: c}
	c.DedicatedVirtualAccounts = &DedicatedVirtualAccountService{client: c}
	c.SplitPayments = &SplitPaymentService{client: c}
	c.PaymentPages = &PaymentPageService{client: c}
	c.BulkCharges = &BulkChargeService{client: c}
	c.Disputes = &DisputeService{client: c}
	c.ApplePayDomains = &ApplePayDomainService{client: c}

	return c, nil
}

// String implements fmt.Stringer so that fmt verbs such as %v and %+v
// never print the client's secret key, even when a caller accidentally
// logs the Client value itself.
func (c *Client) String() string {
	return fmt.Sprintf("paystack.Client{baseURL: %q}", c.baseURL)
}

// GoString implements fmt.GoStringer for the same reason as String, so
// that %#v also omits the secret key.
func (c *Client) GoString() string {
	return c.String()
}
