package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
)

// Respond to requests with an error message
type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenClaims struct {
	jwt.StandardClaims

	// https://openid.net/specs/openid-connect-core-1_0.html#Claims
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	PreferredUsername string `json:"preferred_username"`
	Picture           string `json:"picture"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Nonce             string `json:"nonce"`
}

// Cached signing keypair
var (
	cachedPubKey *rsa.PublicKey
	cachedPrvKey *rsa.PrivateKey
)

// func SigningKeyPair() (pub *rsa.PublicKey, prv *rsa.PrivateKey, err error)
