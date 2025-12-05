package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

func InitJWT() error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("jwtSecret missing")
	}

	jwtKey = []byte(jwtSecret)

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes), err
}

func ComparePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedPassword))

	return err == nil
}

func GenerateJWTToken(userID, fullName, role string) (string, error) {
	if len(jwtKey) == 0 {
		return "", errors.New("missing jwtKey")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   userID,
		"name": fullName,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	if len(jwtKey) == 0 {
		return nil, errors.New("missing jwtKey")
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid sign in method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid Token")
	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid Token format")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claim format")
	}

	return claims, nil
}
