package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type MyClaims struct {
	UserID  string
	IsAdmin bool
	jwt.StandardClaims
}

func GenerateToken(id string, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		UserID:  id,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "Waysbucks",
		},
	})

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, tokenErr := token.SignedString(secretKey)
	if tokenErr != nil {
		return "", tokenErr
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
