package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	verifier "github.com/yukia3e/gcp-sa-jwt/server/jwt"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serviceAccount := os.Getenv("SERVICE_ACCOUNT_EMAIL")

	var jwt string
	flag.StringVar(&jwt, "j", "", "jwt")
	flag.Parse()

	token, err := verifier.Verify(jwt, serviceAccount)
	if err != nil {
		fmt.Errorf("Error: %w", err.Error())
	} else {
		fmt.Println("Success: %+v", token)
	}
}
