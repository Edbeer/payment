package utils

import (
	"os"

	"github.com/Edbeer/auth/types"
	"github.com/golang-jwt/jwt/v4"
)

// Create JWT
func CreateJWT(account *types.Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        account.ID.String(),
		"card":      account.CardNumber,
		"expire_at": 15000,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}