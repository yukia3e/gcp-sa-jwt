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

	aud := os.Getenv("AUD")
	serviceAccount := os.Getenv("SERVICE_ACCOUNT_EMAIL")

	var jwt string
	flag.StringVar(&jwt, "j", "", "jwt")
	flag.Parse()

	fmt.Println(verifier.Verify(jwt, aud, serviceAccount))
}
