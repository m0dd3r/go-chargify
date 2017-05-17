package chargify

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProductsService_List_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"product": {"id":1}},{"product": {"id":2}}]`)
	})

	products, _, err := client.Products.List(context.Background())
	if err != nil {
		t.Errorf("Products.List returned error: %v", err)
	}

	want := []*Product{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(products, want) {
		t.Errorf("Products.List returned %+v, want %+v", products, want)
	}
}
