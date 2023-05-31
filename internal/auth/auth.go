package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("temp")

type JWTClaim struct {
	Code  string   `json:"code"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

func GenerateJWT(email string, code string, roles []string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email: email,
		Code:  code,
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}

	return claims, nil
}
