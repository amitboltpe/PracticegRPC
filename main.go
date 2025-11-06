package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// global variable to store key set
var jwks jwk.Set
var privateKey *rsa.PrivateKey

func init() {

	var err error

	// 1. Generate RSA key pair
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate RSA key: %v", err)
	}

	// 2. Create a JWK (JSON Web Key)
	key, err := jwk.New(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("failed to create JWK: %v", err)
	}

	key.Set("alg", "RS256")
	key.Set("use", "sig")
	key.Set("kid", "Amit-Jain-Key-ID-1234")

	// 3. Create a JWKS (set of keys)
	jwks = jwk.NewSet()
	jwks.Add(key)

	// 4. Serve JWKS endpoint
	http.HandleFunc("./jwks.json", jwksHandler)
	http.HandleFunc("/token", tokenHandler)

	log.Println("JWKS server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))

}

func main() {
	fmt.Println("hellos")

	token := "esfhjfdhjfhjsfjshfbjsfdbdjsffjdshbhvjjbd"
	err := verifyJWT(token)
	if err != nil {
		log.Fatalf(" Token invalid: %v", err)
	}
}

// JWKS endpoint
func jwksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	// 1️⃣ Create a JWT token and set claims
	tok := jwt.New()
	tok.Set("sub", "amit1234")
	tok.Set("iss", "jwks-demo-part")
	tok.Set("role", "admin")
	tok.Set("iat", time.Now().Unix())
	tok.Set("exp", time.Now().Add(5*time.Minute).Unix())

	// 2️⃣ Sign using your RSA private key with RS256 algorithm
	//signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, privateKey))

	claims := tok.PrivateClaims() // map[string]interface{}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		panic(err)
	}
	fmt.Println("JWT claims as bytes:", claimsBytes)

	if err != nil {
		http.Error(w, "failed to sign token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("string(claimsBytes) : ", string(claimsBytes))
	// 3️⃣ Return the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": string(claimsBytes),
	})
	fmt.Println("string(claimsBytes) : ", string(claimsBytes))

}

func verifyJWT(tokenString string) error {
	// Fetch the JWKS from your endpoint
	set, err := jwk.Fetch(context.Background(), "http://localhost:8081/jwks.json")
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	// Parse and verify using the key set
	tok, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(set))
	if err != nil {
		return fmt.Errorf("failed to parse/verify token: %w", err)
	}

	// If you want to check claims:
	sub, _ := tok.Get("sub")
	fmt.Println("✅ JWT verified. Subject:", sub)
	return nil
}
