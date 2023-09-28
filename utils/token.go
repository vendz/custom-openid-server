package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenPayload struct {
	Id   primitive.ObjectID
	Role string
}

func GenerateToken(payload, rand, secretJWTKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = payload
	claims["rand"] = rand
	claims["exp"] = time.Now().AddDate(100, 0, 0).Unix()

	tokenString, err := token.SignedString([]byte(secretJWTKey))

	if err != nil {
		return "", fmt.Errorf("generating JWT Token failed: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(token string, signedJWTKey string) (string, string, error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return []byte(signedJWTKey), nil
	})
	if err != nil {
		return "", "", fmt.Errorf("invalidate token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return "", "", fmt.Errorf("invalid token claim")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid email claim")
	}
	rand, ok := claims["rand"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid issued at claim")
	}

	return email, rand, nil
}
