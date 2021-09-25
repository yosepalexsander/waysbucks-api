package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

type MyClaims struct {
	UserID int 
	jwt.StandardClaims
}

const TokenCtxKey = contextKey("tokenPayload")

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authValue := strings.TrimSpace(r.Header.Get("Authorization"))

		if (len(authValue) == 0) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		bearerSplits := strings.Fields(authValue);

		if (len(bearerSplits) != 2 || bearerSplits[0] != "Bearer") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Authorization header is invalid"))
			return
		}

		token, err := jwt.ParseWithClaims(bearerSplits[1], &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("token is not valid anymore"))
			return
		}
		
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := token.Claims
		ctx := context.WithValue(r.Context(), TokenCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}