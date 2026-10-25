package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func ExtractToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		tokenString = r.URL.Query().Get("token")
	}
	if tokenString == "" {
		return "", errors.New("missing token string")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	return tokenString, nil
}

func GetClaims(r *http.Request) jwt.MapClaims {
	claims, _ := r.Context().Value(ClaimsKey).(jwt.MapClaims)
	return claims
}

func GetUserID(r *http.Request) (string, error) {
	val := r.Context().Value(UserIDKey)
	if val == nil {
		return "", errors.New("userID not found in context")
	}
	userID, ok := val.(string)
	if !ok {
		return "", errors.New("invalid userID in context")
	}
	return userID, nil
}
