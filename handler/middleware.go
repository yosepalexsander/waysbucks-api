package handler

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

type MyClaims struct {
	UserID uint64 
	jwt.StandardClaims
}

const tokenCtxKey = contextKey("tokenPayload")

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authValue := strings.TrimSpace(r.Header.Get("Authorization"))

		if (len(authValue) == 0) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		bearerSplits := strings.Fields(authValue);

		if (len(bearerSplits) != 2 || bearerSplits[0] != "Bearer") {
			badRequest(w, "Authorization header value is invalid")
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
			badRequest(w, "token is invalid anymore")
			return
		}
		
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*MyClaims)
		ctx := context.WithValue(r.Context(), tokenCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}