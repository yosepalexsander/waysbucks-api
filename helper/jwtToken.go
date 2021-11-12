package helper

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yosepalexsander/waysbucks-api/entity"
)

type MyClaims struct {
	UserID  int
	IsAdmin bool
	jwt.StandardClaims
}

func GenerateToken(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		UserID:  user.Id,
		IsAdmin: user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "Waysbucks",
		},
	})

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, tokenErr := token.SignedString(secretKey)
	if tokenErr != nil {
		log.Println(tokenErr)
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
