package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/yosepalexsander/waysbucks-api/helper"
)

type contextKey struct {
	name string
}

var TokenCtxKey = &contextKey{name: "tokenPayload"}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authValue := strings.Fields(strings.TrimSpace(r.Header.Get("Authorization")))

		if len(authValue) != 2 {
			http.Error(w, "Authorization header is invalid", http.StatusBadRequest)
			return
		}

		if authValue[0] != "Bearer" {
			http.Error(w, "Authorization header is invalid", http.StatusBadRequest)
			return
		}

		token, err := helper.VerifyToken(authValue[1])

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "token is not valid anymore", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "token is not valid anymore", http.StatusUnauthorized)
			return
		}

		claims := token.Claims
		ctx := context.WithValue(r.Context(), TokenCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(TokenCtxKey).(*helper.MyClaims)
		if !ok || !claims.IsAdmin {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
