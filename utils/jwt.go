package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

func GenrateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	},
	)
	return token.SignedString([]byte(secretKey))
}

// parse the JWT token and validate its signature and method.
func parseAndValidateToken(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("error parsing token")
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return parsedToken, nil
}

// extracts claims from a validated JWT token.
func extractClaims(parsedToken *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

// Attempt to extract value from claim, check it matches the expected type anc return it.
func extractFromClaim[T any](claims jwt.MapClaims, key string) (T, error) {
	value, ok := claims[key]
	if !ok {
		return *new(T), errors.New(key + " not found in token claims")
	}

	// Perform type assertion to convert value to the specified type
	typedValue, ok := value.(T)
	if !ok {
		return *new(T), errors.New(key + " is not of expected type")
	}

	return typedValue, nil
}

func VerifyToken(token string) (int64, error) {
	// Parse and validate the token
	parsedToken, err := parseAndValidateToken(token)
	if err != nil {
		return 0, err
	}

	// Extract claims from the token
	claims, err := extractClaims(parsedToken)
	if err != nil {
		return 0, err
	}

	// Extract userId from the claims
	userId, err := extractFromClaim[float64](claims, "userId")
	if err != nil {
		return 0, err
	}

	return int64(userId), nil
}
