package tests

import "testing"

func TestListProducts(t *testing.T) {
	_, _, err := client.Products.List(ctx)
	if err != nil {
		t.Error(err)
	}

}
