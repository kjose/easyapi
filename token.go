// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package easyapi

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/kjose/jgmc/api/internal/easyapi/layer"
)

type Token struct {
	Value     string
	ExpiresAt int
}

type claims struct {
	Info interface{}
	jwt.StandardClaims
}

// Generate a JWT token
func GenerateToken(t layer.TokenInterface, duration time.Duration) *Token {
	jwtKey := []byte(os.Getenv("JWT_TOKEN_KEY"))
	expirationDate := time.Now().Add(duration).Unix()
	claims := claims{
		Info: t.GetTokenInformations(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationDate,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return &Token{
		Value:     tokenString,
		ExpiresAt: int(expirationDate),
	}
}

// Parse a string into a valid JWT token and returns the token information or return an error
func ParseToken(tokenString string) (interface{}, error) {
	jwtKey := []byte(os.Getenv("JWT_TOKEN_KEY"))
	var claims claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return claims.Info, nil
}
