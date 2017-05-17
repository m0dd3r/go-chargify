package chargify

import (
	"context"
	"fmt"
)

type SubscriptionWrapper struct {
	Subscription *Subscription `json:"subscription"`
}

type Subscription struct {
	Id                          int            `json:"id,omitempty"`
	State                       string         `json:"state,omitempty"`
	TrialStartedAt              *FormattedTime `json:"trial_started_at,omitempty"`
	Customer                    *Customer      `json:"customer,omitempty"`
	CustomerAttributes          *Customer      `json:"customer_attributes,omitempty"`
	Product                     *Product       `json:"product,omitempty"`
	ProductHandle               string         `json:"product_handle,omitempty"`
	CreditCard                  *CreditCard    `json:"credit_card,omitempty"`
	TrialEndedAt                *FormattedTime `json:"trial_ended_at,omitempty"`
	ActivatedAt                 *FormattedTime `json:"activated_at,omitempty"`
	CreatedAt                   *FormattedTime `json:"created_at,omitempty"`
	UpdatedAt                   *FormattedTime `json:"updated_at,omitempty"`
	ExpiresAt                   *FormattedTime `json:"expires_at,omitempty"`
	PreviousExpiresAt           *FormattedTime `json:"previous_expires_at,omitempty"`
	BalanceInCents              int            `json:"balance_in_cents,omitempty"`
	CurrentPeriodEndsAt         *FormattedTime `json:"current_period_ends_at,omitempty"`
	NextAssessmentAt            *FormattedTime `json:"next_assessment_at,omitempty"`
	CanceledAt                  *FormattedTime `json:"canceled_at,omitempty"`
	CancellationMessage         string         `json:"cancellation_message,omitempty"`
	NextProductId               int            `json:"next_product_id,omitempty"`
	CancelAtEndOfPeriod         bool           `json:"cancel_at_end_of_period,omitempty"`
	PaymentCollectionMethod     string         `json:"payment_collection_method,omitempty"`
	SnapDay                     string         `json:"snap_day,omitempty"`
	CancellationMethod          string         `json:"cancellation_method,omitempty"`
	CurrentPeriodStartedAt      *FormattedTime `json:"current_period_started_at,omitempty"`
	PreviousState               string         `json:"previous_state,omitempty"`
	SignupPaymentId             int            `json:"signup_payment_id,omitempty"`
	SignupRevenue               float32        `json:"signup_revenue,omitempty,string"`
	DelayedCancelAt             *FormattedTime `json:"delayed_cancel_at,omitempty"`
	CouponCode                  string         `json:"coupon_code,omitempty"`
	TotalRevenueInCents         int            `json:"total_revenue_in_cents,omitempty"`
	ProductPriceInCents         int            `json:"product_price_in_cents,omitempty"`
	ProductVersionNumber        int            `json:"product_version_number,omitempty"`
	PaymentType                 string         `json:"payment_type,omitempty"`
	ReferralCode                string         `json:"referral_code,omitempty"`
	CouponUseCount              int            `json:"coupon_use_count,omitempty"`
	CouponUsesAllowed           int            `json:"coupon_uses_allowed,omitempty"`
	CurrentBillingAmountInCents int            `json:"current_billing_amount_in_cents,omitempty"`
}

type SubscriptionsService service

func (svc *SubscriptionsService) Create(ctx context.Context, sub *Subscription) (*Subscription, *Response, error) {
	u := "/subscriptions"

	sw := SubscriptionWrapper{sub}
	req, err := svc.client.NewRequest("POST", u, sw)
	if err != nil {
		return nil, nil, err
	}

	swr := new(SubscriptionWrapper)
	resp, err := svc.client.Do(ctx, req, swr)
	if err != nil {
		return nil, resp, err
	}

	return swr.Subscription, resp, nil
}

// Get fetches a subscription.
//
// Chargify API docs: https://reference.chargify.com/v1/subscriptions/read-subscription
func (s *SubscriptionsService) Get(ctx context.Context, id int) (*Subscription, *Response, error) {
	u := fmt.Sprintf("subscriptions/%d", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	sw := new(SubscriptionWrapper)
	resp, err := s.client.Do(ctx, req, sw)
	if err != nil {
		return nil, resp, err
	}

	return sw.Subscription, resp, nil
}
