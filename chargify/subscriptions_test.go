package chargify

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

const (
	subscriptionJSON = `{
		"subscription": {
			"id": 14900541,
			"state": "{{.State}}",
			"trial_started_at": "2016-10-24T16:20:12-04:00",
			"trial_ended_at": "2016-10-24T16:20:43-04:00",
			"activated_at": "2016-10-24T16:20:43-04:00",
			"created_at": "2016-10-24T16:20:12-04:00",
			"updated_at": "2016-11-03T09:34:37-04:00",
			"expires_at": null,
			"balance_in_cents": 2450,
			"current_period_ends_at": "2016-12-01T11:41:25-05:00",
			"next_assessment_at": "2016-12-01T11:41:25-05:00",
			"canceled_at": null,
			"cancellation_message": null,
			"next_product_id": null,
			"cancel_at_end_of_period": false,
			"payment_collection_method": "invoice",
			"snap_day": null,
			"cancellation_method": null,
			"current_period_started_at": "2016-11-01T12:41:25-04:00",
			"previous_state": "active",
			"signup_payment_id": 159423810,
			"signup_revenue": "0.00",
			"delayed_cancel_at": null,
			"coupon_code": null,
			"total_revenue_in_cents": 18000,
			"product_price_in_cents": 4000,
			"product_version_number": 4,
			"payment_type": "credit_card",
			"referral_code": "p8fs35",
			"coupon_use_count": null,
			"coupon_uses_allowed": null,
			"current_billing_amount_in_cents": 6450,
			"customer": {
				"id": 14399371,
				"first_name": "Amelia",
				"last_name": "Example",
				"organization": "Acme",
				"email": "amelia@example.com",
				"created_at": "2016-10-24T16:20:12-04:00",
				"updated_at": "2016-10-26T13:25:33-04:00",
				"reference": "JQPUBLIC",
				"address": "123 Anywhere Street",
				"address_2": "",
				"city": "Anywhere",
				"state": "MA",
				"zip": "02120",
				"country": "US",
				"phone": "555-555-1212",
				"portal_invite_last_sent_at": null,
				"portal_invite_last_accepted_at": null,
				"verified": false,
				"portal_customer_created_at": null,
				"cc_emails": "john@example.com, joe@example.com"
			},
			"product": {
				"id": 3792003,
				"name": "$10 Basic Plan",
				"handle": "basic",
				"description": "lorem ipsum",
				"accounting_code": "basic",
				"request_credit_card": false,
				"expiration_interval": null,
				"expiration_interval_unit": "never",
				"created_at": "2016-03-24T13:38:39-04:00",
				"updated_at": "2016-11-03T13:03:05-04:00",
				"price_in_cents": 1000,
				"interval": 1,
				"interval_unit": "day",
				"initial_charge_in_cents": null,
				"trial_price_in_cents": null,
				"trial_interval": null,
				"trial_interval_unit": "month",
				"archived_at": null,
				"require_credit_card": false,
				"return_params": "",
				"taxable": false,
				"update_return_url": "",
				"initial_charge_after_trial": false,
				"version_number": 7,
				"update_return_params": "",
				"product_family": {
					"id": 527890,
					"name": "Acme Projects",
					"description": "",
					"handle": "billing-plans",
					"accounting_code": null
				},
				"public_signup_pages": [
				{
					"id": 281054,
					"return_url": "http://www.example.com?successfulsignup",
					"return_params": "",
					"url": "https://general-goods.chargify.com/subscribe/kqvmfrbgd89q/basic"
				},
				{
					"id": 281240,
					"return_url": "",
					"return_params": "",
					"url": "https://general-goods.chargify.com/subscribe/dkffht5dxfd8/basic"
				},
				{
					"id": 282694,
					"return_url": "",
					"return_params": "",
					"url": "https://general-goods.chargify.com/subscribe/jwffwgdd95s8/basic"
				}
				]
			},
			"credit_card": {
				"id": 9979580,
				"first_name": "Amelia",
				"last_name": "Example",
				"masked_card_number": "XXXX-XXXX-XXXX-1",
				"card_type": "bogus",
				"expiration_month": 1,
				"expiration_year": 2026,
				"customer_id": 14399371,
				"current_vault": "bogus",
				"vault_token": "1",
				"billing_address": "123 Anywhere Street",
				"billing_city": "Anywhere",
				"billing_state": "MA",
				"billing_zip": "02120",
				"billing_country": "US",
				"customer_vault_token": null,
				"billing_address_2": "",
				"payment_type": "credit_card"
			}
		}
	}`
)

func testSubJSON(state string) string {
	t, err := template.New("subscriptionJSON").Parse(subscriptionJSON)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, struct{ State string }{state})
	if err != nil {
		panic(err)
	}
	return buf.String()

}

func testSub(state string) *Subscription {
	if state == "" {
		state = "active"
	}
	return &Subscription{
		Id:                          14900541,
		State:                       state,
		TrialStartedAt:              NewFormattedTime(`"2016-10-24T16:20:12-04:00"`),
		TrialEndedAt:                NewFormattedTime(`"2016-10-24T16:20:43-04:00"`),
		ActivatedAt:                 NewFormattedTime(`"2016-10-24T16:20:43-04:00"`),
		CreatedAt:                   NewFormattedTime(`"2016-10-24T16:20:12-04:00"`),
		UpdatedAt:                   NewFormattedTime(`"2016-11-03T09:34:37-04:00"`),
		ExpiresAt:                   nil,
		BalanceInCents:              2450,
		CurrentPeriodEndsAt:         NewFormattedTime(`"2016-12-01T11:41:25-05:00"`),
		NextAssessmentAt:            NewFormattedTime(`"2016-12-01T11:41:25-05:00"`),
		CanceledAt:                  nil,
		CancellationMessage:         "",
		NextProductId:               0,
		CancelAtEndOfPeriod:         false,
		PaymentCollectionMethod:     "invoice",
		SnapDay:                     "",
		CancellationMethod:          "",
		CurrentPeriodStartedAt:      NewFormattedTime(`"2016-11-01T12:41:25-04:00"`),
		PreviousState:               "active",
		SignupPaymentId:             159423810,
		SignupRevenue:               0.00,
		DelayedCancelAt:             nil,
		CouponCode:                  "",
		TotalRevenueInCents:         18000,
		ProductPriceInCents:         4000,
		ProductVersionNumber:        4,
		PaymentType:                 "credit_card",
		ReferralCode:                "p8fs35",
		CouponUseCount:              0,
		CouponUsesAllowed:           0,
		CurrentBillingAmountInCents: 6450,
		Customer: &Customer{
			Id:           14399371,
			FirstName:    "Amelia",
			LastName:     "Example",
			Organization: "Acme",
			Email:        "amelia@example.com",
			CreatedAt:    NewFormattedTime(`"2016-10-24T16:20:12-04:00"`),
			UpdatedAt:    NewFormattedTime(`"2016-10-26T13:25:33-04:00"`),
			Reference:    "JQPUBLIC",
			Address:      "123 Anywhere Street",
			Address2:     "",
			City:         "Anywhere",
			State:        "MA",
			Zip:          "02120",
			Country:      "US",
			Phone:        "555-555-1212",
			PortalInviteLastSentAt:     nil,
			PortalInviteLastAcceptedAt: nil,
			Verified:                   false,
			PortalCustomerCreatedAt:    nil,
			CcEmails:                   "john@example.com, joe@example.com",
		},
		Product: &Product{
			Id:                      3792003,
			Name:                    "$10 Basic Plan",
			Handle:                  "basic",
			Description:             "lorem ipsum",
			AccountingCode:          "basic",
			RequestCreditCard:       false,
			ExpirationInterval:      0,
			ExpirationIntervalUnit:  "never",
			CreatedAt:               NewFormattedTime(`"2016-03-24T13:38:39-04:00"`),
			UpdatedAt:               NewFormattedTime(`"2016-11-03T13:03:05-04:00"`),
			PriceInCents:            1000,
			Interval:                1,
			IntervalUnit:            "day",
			InitialChargeInCents:    0,
			TrialPriceInCents:       0,
			TrialInterval:           0,
			TrialIntervalUnit:       "month",
			ArchivedAt:              nil,
			RequireCreditCard:       false,
			ReturnParams:            "",
			Taxable:                 false,
			UpdateReturnUrl:         "",
			InitialChargeAfterTrial: false,
			VersionNumber:           7,
			UpdateReturnParams:      "",
			ProductFamily: &ProductFamily{
				Id:             527890,
				Name:           "Acme Projects",
				Description:    "",
				Handle:         "billing-plans",
				AccountingCode: "",
			},
			PublicSignupPages: []*PublicSignupPage{
				&PublicSignupPage{
					Id:           281054,
					ReturnUrl:    "http://www.example.com?successfulsignup",
					ReturnParams: "",
					Url:          "https://general-goods.chargify.com/subscribe/kqvmfrbgd89q/basic",
				},
				&PublicSignupPage{
					Id:           281240,
					ReturnUrl:    "",
					ReturnParams: "",
					Url:          "https://general-goods.chargify.com/subscribe/dkffht5dxfd8/basic",
				},
				&PublicSignupPage{
					Id:           282694,
					ReturnUrl:    "",
					ReturnParams: "",
					Url:          "https://general-goods.chargify.com/subscribe/jwffwgdd95s8/basic",
				},
			},
		},
		CreditCard: &CreditCard{
			Id:                 9979580,
			FirstName:          "Amelia",
			LastName:           "Example",
			MaskedCardNumber:   "XXXX-XXXX-XXXX-1",
			CardType:           "bogus",
			ExpirationMonth:    1,
			ExpirationYear:     2026,
			CustomerId:         14399371,
			CurrentVault:       "bogus",
			VaultToken:         "1",
			BillingAddress:     "123 Anywhere Street",
			BillingCity:        "Anywhere",
			BillingState:       "MA",
			BillingZip:         "02120",
			BillingCountry:     "US",
			CustomerVaultToken: "",
			BillingAddress2:    "",
			PaymentType:        "credit_card",
		},
	}
}

func TestSubscriptionsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/subscriptions/14900541", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, testSubJSON("active"))
	})

	sub, _, err := client.Subscriptions.Get(context.Background(), 14900541)
	if err != nil {
		t.Errorf("Subscription.Get returned error: %v", err)
	}

	want := testSub("active")
	if diff := pretty.Compare(sub, want); diff != "" {
		t.Errorf("Products.List diff: (-got +want)\n%s\n", diff)
	}
}

func TestSubscriptionsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, testSubJSON("active"))
	})

	input := testSub("active")
	input.Id = 0
	sub, _, err := client.Subscriptions.Create(context.Background(), input)
	if err != nil {
		t.Errorf("Subscription.Get returned error: %v", err)
	}

	want := testSub("active")
	if diff := pretty.Compare(sub, want); diff != "" {
		t.Errorf("Products.List diff: (-got +want)\n%s\n", diff)
	}
}

func TestSubscriptionsService_Destroy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/subscriptions/14900541", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, testSubJSON("canceled"))
	})

	sub, _, err := client.Subscriptions.Destroy(context.Background(), 14900541)
	if err != nil {
		t.Errorf("Subscription.Destroy returned error: %v", err)
	}

	want := testSub("canceled")
	if diff := pretty.Compare(sub, want); diff != "" {
		t.Errorf("Products.List diff: (-got +want)\n%s\n", diff)
	}
}
