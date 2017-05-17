package chargify

import "context"

type ProductWrapper struct {
	Product *Product `json:"product"`
}

type Product struct {
	Id                      int                 `json:"id,omitempty"`
	Name                    string              `json:"name,omitempty"`
	Handle                  string              `json:"handle,omitempty"`
	Description             string              `json:"description,omitempty"`
	AccountingCode          string              `json:"accounting_code,omitempty"`
	RequestCreditCard       bool                `json:"request_credit_card,omitempty"`
	ExpirationInterval      int                 `json:"expiration_interval,omitempty"`
	ExpirationIntervalUnit  string              `json:"expiration_interval_unit,omitempty"`
	CreatedAt               *FormattedTime      `json:"created_at,omitempty"`
	UpdatedAt               *FormattedTime      `json:"updated_at,omitempty"`
	PriceInCents            int                 `json:"price_in_cents,omitempty"`
	Interval                int                 `json:"interval,omitempty"`
	IntervalUnit            string              `json:"interval_unit,omitempty"`
	InitialChargeInCents    int                 `json:"initial_charge_in_cents,omitempty"`
	TrialPriceInCents       int                 `json:"trial_price_in_cents,omitempty"`
	TrialInterval           int                 `json:"trial_interval,omitempty"`
	TrialIntervalUnit       string              `json:"trial_interval_unit,omitempty"`
	ArchivedAt              *FormattedTime      `json:"archived_at,omitempty"`
	RequireCreditCard       bool                `json:"require_credit_card,omitempty"`
	ReturnParams            string              `json:"return_params,omitempty"`
	Taxable                 bool                `json:"taxable,omitempty"`
	UpdateReturnUrl         string              `json:"update_return_url,omitempty"`
	InitialChargeAfterTrial bool                `json:"initial_charge_after_trial,omitempty"`
	VersionNumber           int                 `json:"version_number,omitempty"`
	UpdateReturnParams      string              `json:"update_return_params,omitempty"`
	ProductFamily           *ProductFamily      `json:"product_family,omitempty"`
	PublicSignupPages       []*PublicSignupPage `json:"public_signup_pages,omitempty"`
}

type ProductFamily struct {
	Id             int    `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Handle         string `json:"handle,omitempty"`
	Description    string `json:"description,omitempty"`
	AccountingCode string `json:"accounting_code,omitempty"`
}

type PublicSignupPage struct {
	Id           int    `json:"id,omitempty"`
	ReturnUrl    string `json:"return_url,omitempty"`
	ReturnParams string `json:"return_params,omitempty"`
	Url          string `json:"url,omitempty"`
}

type ProductsService service

// List fetches all products.
//
// Chargify API docs: https://reference.chargify.com/v1/products/list-products
func (s *ProductsService) List(ctx context.Context) ([]*Product, *Response, error) {
	u := "products"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var wrappers []*ProductWrapper
	resp, err := s.client.Do(ctx, req, &wrappers)
	if err != nil {
		return nil, resp, err
	}
	var products []*Product
	for _, p := range wrappers {
		products = append(products, p.Product)
	}
	return products, resp, nil
}
