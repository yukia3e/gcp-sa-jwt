package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"github.com/joho/godotenv"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

type claimSet struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	aud := os.Getenv("AUD")
	serviceAccount := os.Getenv("SERVICE_ACCOUNT_EMAIL")

	log.Print("Hello world received a request.")
	ctx := context.Background()
	c, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		// TODO: Handle error.
		fmt.Fprintf(w, "Error:%s", err.Error())
		return
	}
	defer c.Close()

	var jwt claimSet
	serviceAccountName := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccount)
	oauthTokenURI := aud

	jwt.Iss = serviceAccountName
	jwt.Sub = serviceAccountName
	jwt.Aud = oauthTokenURI
	jwt.Iat = time.Now().Unix()
	jwt.Exp = time.Now().Add(time.Hour).Unix()

	claimSetJSON, err := json.Marshal(jwt)
	if err != nil {
		// TODO: Handle error.
		fmt.Fprintf(w, "Error:%s", err.Error())
		return
	}
	// fmt.Fprintf(w, "claimSetJSON:%s", claimSetJSON)

	// See https://pkg.go.dev/google.golang.org/genproto/googleapis/iam/credentials/v1#SignJwtRequest.
	req := &credentialspb.SignJwtRequest{
		Name: serviceAccountName,
		// See https://developers.google.com/identity/protocols/oauth2/service-account#jwt-auth
		Payload: string(claimSetJSON),
	}
	resp, err := c.SignJwt(ctx, req)
	if err != nil {
		// TODO: Handle error.
		fmt.Fprintf(w, "Error:%s", err.Error())
		return
	}
	fmt.Fprintf(w, "JWT %s!\n", resp)
}
