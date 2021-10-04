package middleware

import (
	"net/http"

	"github.com/yosepalexsander/waysbucks-api/helper"
)
 
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