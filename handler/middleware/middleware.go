package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/yosepalexsander/waysbucks-api/helper"
)

type contextKey struct {
	name string
}

var TokenCtxKey = &contextKey{name: "tokenPayload"}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")
		if err != nil || c.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := helper.VerifyToken(c.Value)

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