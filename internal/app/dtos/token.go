package dtos

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	Username string `json:"username,omitempty"`
	jwt.RegisteredClaims
}
