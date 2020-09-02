package model

import "github.com/dgrijalva/jwt-go"

// JWTClaims is a truct that will be encoded to a JWT and add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type JWTClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
