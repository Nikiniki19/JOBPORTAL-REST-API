package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const Key ctxKey = 1

// Auth Struct
type Auth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}
//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=auth
type Authentication interface {
	GenerateToken(claims jwt.RegisteredClaims) (string, error)
	ValidateToken(token string) (jwt.RegisteredClaims, error)
}

// Creating NewAuth Factory Function
func NewAuth(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (Authentication, error) {
	if privateKey == nil || publicKey == nil {
		return nil, errors.New("private/public key cannot be nil")
	}
	return &Auth{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// Generating Tokens
func (a *Auth) GenerateToken(claims jwt.RegisteredClaims) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token %w", err)
	}
	return tokenStr, nil
}

// Validating the tokens
func (a *Auth) ValidateToken(token string) (jwt.RegisteredClaims, error) {
	var c jwt.RegisteredClaims
	tkn, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("parsing token %w", err)
	}
	if !tkn.Valid {
		return jwt.RegisteredClaims{}, errors.New("invalid token")
	}
	return c, nil
}
