package jwtparser

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

// User Custom object which can be stored in the claims
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// AuthTokenClaim This is the claim object which gets parsed from the authorization header
type AuthTokenClaim struct {
	*jwt.StandardClaims
	User
}

//TokenManager is an interface for methods to extract token related information
type TokenManager interface {
	ExtractJWTClaims(bearerToken string) (*AuthTokenClaim, error)
}

//JWTTokenManager is a concrete class which implements TokenManager interface
type JWTTokenManager struct {
}

//NewJWTTokenManager is a constructor to create JWTTokenManager object
func NewJWTTokenManager() *JWTTokenManager {
	return &JWTTokenManager{}
}

//ExtractJWTClaims extracts custom claims from a valid token
func (t *JWTTokenManager) ExtractJWTClaims(bearerToken string) (*AuthTokenClaim, error) {
	userClaims := &AuthTokenClaim{}

	tokenString := strings.Split(bearerToken, " ")

	token, err := jwt.ParseWithClaims(tokenString[1], userClaims, func(token *jwt.Token) (interface{}, error) {
		secret := os.Getenv("JWT_TOKEN_SECRET")
		if len(secret) == 0 {
			log.Printf("error: JWT_TOKEN_SECRET env variable is not set")
			return nil, errors.New("authentication secret not found")
		}
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthTokenClaim); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token claim is not ok")
}
