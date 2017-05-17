package tests

import (
	"testing"

	"github.com/m0dd3r/go-chargify/chargify"
)

func TestCreateSubscription(t *testing.T) {
	_, _, err := client.Subscriptions.Create(ctx, &chargify.Subscription{
		CustomerAttributes: &chargify.Customer{
			FirstName: "Bob",
			LastName:  "Test",
			Email:     "foo@example.com",
		},
		ProductHandle: "standard",
	})
	if err != nil {
		t.Error(err)
	}
}
