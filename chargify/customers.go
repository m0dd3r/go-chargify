package chargify

type Customer struct {
	Id                         int            `json:"id,omitempty"`
	FirstName                  string         `json:"first_name,omitempty"`
	LastName                   string         `json:"last_name,omitempty"`
	Organization               string         `json:"organization,omitempty"`
	Email                      string         `json:"email,omitempty"`
	CreatedAt                  *FormattedTime `json:"created_at,omitempty"`
	UpdatedAt                  *FormattedTime `json:"updated_at,omitempty"`
	Reference                  string         `json:"reference,omitempty"`
	Address                    string         `json:"address,omitempty"`
	Address2                   string         `json:"address_2,omitempty"`
	City                       string         `json:"city,omitempty"`
	State                      string         `json:"state,omitempty"`
	Zip                        string         `json:"zip,omitempty"`
	Country                    string         `json:"country,omitempty"`
	Phone                      string         `json:"phone,omitempty"`
	PortalInviteLastSentAt     *FormattedTime `json:"portal_invite_last_sent_at,omitempty"`
	PortalInviteLastAcceptedAt *FormattedTime `json:"portal_invite_last_accepted_at,omitempty"`
	Verified                   bool           `json:"verified,omitempty"`
	PortalCustomerCreatedAt    *FormattedTime `json:"portal_customer_created_at,omitempty"`
	CcEmails                   string         `json:"cc_emails,omitempty"`
	TaxExempt                  bool           `json:"tax_exempt,omitempty"`
}

type CreditCard struct {
	Id                 int    `json:"id,omitempty"`
	FirstName          string `json:"first_name,omitempty"`
	LastName           string `json:"last_name,omitempty"`
	MaskedCardNumber   string `json:"masked_card_number,omitempty"`
	CardType           string `json:"card_type,omitempty"`
	ExpirationMonth    int    `json:"expiration_month,omitempty"`
	ExpirationYear     int    `json:"expiration_year,omitempty"`
	CustomerId         int    `json:"customer_id,omitempty"`
	CurrentVault       string `json:"current_vault,omitempty"`
	VaultToken         string `json:"vault_token,omitempty"`
	BillingAddress     string `json:"billing_address,omitempty"`
	BillingCity        string `json:"billing_city,omitempty"`
	BillingState       string `json:"billing_state,omitempty"`
	BillingZip         string `json:"billing_zip,omitempty"`
	BillingCountry     string `json:"billing_country,omitempty"`
	CustomerVaultToken string `json:"customer_vault_token,omitempty"`
	BillingAddress2    string `json:"billing_address_2,omitempty"`
	PaymentType        string `json:"payment_type,omitempty"`
}
