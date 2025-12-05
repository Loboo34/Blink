package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func ExtractToken(r *http.Request)(string, error){
	tokenString := r.Header.Get("Authorization")
	if tokenString == ""{
		return "", errors.New("missing token string")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	return tokenString, nil
}

func GetClaims(r *http.Request)jwt.MapClaims{
	claims,_ := r.Context().Value("claims").(jwt.MapClaims)
	return claims
}

func GetUserID(r *http.Request)(string, error){
userID := r.Context().Value("userID").(string)
return userID, nil
}

