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
		authValue := strings.TrimSpace(r.Header.Get("Authorization"))

		if (len(authValue) == 0) {
			http.Error(w, "authorization header is invalid", http.StatusBadRequest)
			return
		}

		bearerSplits := strings.Fields(authValue);

		if (len(bearerSplits) != 2 || bearerSplits[0] != "Bearer") {
			http.Error(w, "authorization header is invalid", http.StatusBadRequest)
			return
		}

		token, err := helper.VerifyToken(bearerSplits[1])

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "token is not valid anymore", http.StatusBadRequest)
			return
		}
		
		if !token.Valid {
			http.Error(w, "token is not valid anymore", http.StatusBadRequest)
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