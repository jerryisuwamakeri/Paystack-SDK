package paystack

import (
	"context"
	"net/url"
	"strconv"
)

// ProductService groups the product-related Paystack API operations.
// Access it via Client.Products.
type ProductService struct {
	client *Client
}

// Product represents a Paystack payment page product.
type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Currency    string `json:"currency"`
	ProductCode string `json:"product_code"`
	Quantity    int64  `json:"quantity"`
	Unlimited   bool   `json:"unlimited"`
}

// CreateProductRequest is the payload for ProductService.Create.
type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Currency    string `json:"currency,omitempty"`
	Unlimited   bool   `json:"unlimited,omitempty"`
	Quantity    int64  `json:"quantity,omitempty"`
}

// Create creates a new product.
func (s *ProductService) Create(ctx context.Context, req CreateProductRequest, opts ...RequestOption) (*Product, *Response, error) {
	if req.Name == "" {
		return nil, nil, ErrMissingName
	}
	if req.Price <= 0 {
		return nil, nil, ErrInvalidAmount
	}

	var out Product
	resp, err := s.client.do(ctx, "POST", "/product", req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Fetch retrieves a product by its numeric ID.
func (s *ProductService) Fetch(ctx context.Context, id string, opts ...RequestOption) (*Product, *Response, error) {
	if id == "" {
		return nil, nil, ErrMissingProductID
	}

	var out Product
	resp, err := s.client.do(ctx, "GET", "/product/"+url.PathEscape(id), nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// UpdateProductRequest is the payload for ProductService.Update.
type UpdateProductRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price,omitempty"`
	Currency    string `json:"currency,omitempty"`
}

// Update modifies an existing product.
func (s *ProductService) Update(ctx context.Context, id string, req UpdateProductRequest, opts ...RequestOption) (*Product, *Response, error) {
	if id == "" {
		return nil, nil, ErrMissingProductID
	}

	var out Product
	resp, err := s.client.do(ctx, "PUT", "/product/"+url.PathEscape(id), req, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// ListProductsParams filters ProductService.List and ProductService.ListAll.
type ListProductsParams struct {
	Page    int
	PerPage int
}

func (p ListProductsParams) query() url.Values {
	q := url.Values{}
	if p.PerPage > 0 {
		q.Set("perPage", strconv.Itoa(p.PerPage))
	}
	if p.Page > 0 {
		q.Set("page", strconv.Itoa(p.Page))
	}
	return q
}

// List retrieves a single page of products.
func (s *ProductService) List(ctx context.Context, params ListProductsParams, opts ...RequestOption) ([]Product, *Response, error) {
	path := "/product"
	if enc := params.query().Encode(); enc != "" {
		path += "?" + enc
	}

	var out []Product
	resp, err := s.client.do(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// ListAll returns an Iterator over every product, fetching additional
// pages on demand as the caller advances it.
func (s *ProductService) ListAll(ctx context.Context, params ListProductsParams, opts ...RequestOption) *Iterator[Product] {
	return NewIterator(ctx, params.PerPage, func(ctx context.Context, page, perPage int) ([]Product, *Response, error) {
		p := params
		p.Page = page
		p.PerPage = perPage
		return s.List(ctx, p, opts...)
	})
}
