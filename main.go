package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	CLIENT_ID     = GetEnv("CLIENT_ID")
	CLIENT_SECRET = GetEnv("CLIENT_SECRET")
)

var ctx = context.TODO()

const REDIRECT_URL = "https://localhost:8080"

func GetEnv(variable string) string {
	env := os.Getenv(variable)
	if env == "" {
		log.Fatalf("Environment variable %s is not set", variable)
	}
	return env
}

func main() {
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	conf := &oauth2.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		RedirectURL:  "YOUR_REDIRECT_URL",
		Scopes: []string{
			"https://www.googleapis.com/auth/bigquery",
			"https://www.googleapis.com/auth/blogger",
		},
		Endpoint: google.Endpoint,
	}
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state")
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Handle the exchange code to initiate a transport.
	tok, err := conf.Exchange(ctx, "authorization-code")
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(ctx, tok)
	client.Get("...")
}
