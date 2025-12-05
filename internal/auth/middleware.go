package auth

import (
	"context"
	"net/http"

	"github.com/Loboo34/Blink/internal/utils"
)

// ctxKey is an unexported type for context keys defined in this package.
type ctxKey string

const (
	claimsKey ctxKey = "claims"
	userIDKey ctxKey = "userID"
)

func CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := ExtractToken(r)
		if err != nil {
			utils.Logger.Warn("Failed to extract token")
			utils.RespondWithError(w, http.StatusBadRequest, "Missing JWT token", "")
			return
		}

		claims, err := ValidateToken(token)
		if err != nil {
			utils.Logger.Warn("Failed to validate token")
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Token", "")
			return
		}

		userID, ok := claims["id"].(string)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid User ID", "")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		ctx = context.WithValue(ctx, userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
