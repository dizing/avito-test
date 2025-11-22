package domain

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	CLAIMS_USERNAME     = "Username"
	CLAIMS_EXP          = "exp"
	SECRET              = "secret"
	CONTEXT_CREDENTIALS = "credentials"
)

func CreateJwt(username UserName) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		CLAIMS_USERNAME: username,
		CLAIMS_EXP:      time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(SECRET))

	return tokenString
}

func ExtractClaimsFromJwt(tokenString string) (UserName, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("authorization token invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("can't parse token claims")
	}

	username, ok := claims[CLAIMS_USERNAME]
	if !ok {
		return "", fmt.Errorf("claims must contain")
	}

	return UserName(username.(string)), nil
}
