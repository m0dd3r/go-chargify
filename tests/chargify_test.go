package tests

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/m0dd3r/go-chargify/chargify"
)

var (
	client *chargify.Client
	ctx    context.Context

	auth bool
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	subdomain := os.Getenv("CHARGIFY_SUBDOMAIN")
	if subdomain == "" {
		panic("CAN'T RUN TESTS WITHOUT SUBDOMAIN")
	}

	api_key := os.Getenv("CHARGIFY_API_KEY")
	if api_key == "" {
		fmt.Printf("NO API KEY, SOME TESTS MAY NOT RUN")
	} else {
		auth = true
	}
	ctx = context.Background()
	client = chargify.NewClient(subdomain, api_key, nil)
	fmt.Println("TEST CLIENT INITIALIZED")
}
