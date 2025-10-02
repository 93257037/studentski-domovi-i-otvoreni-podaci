package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JWTClaims - podaci koji se cuvaju u JWT tokenu
type JWTClaims struct {
	UserID   primitive.ObjectID `json:"user_id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Role     string             `json:"role"`
	jwt.RegisteredClaims
}

// generiše JWT token za korisnika sa prosledjenim podacima
// token vazi 24 sata i potpisan je sa tajnim kljucem
func GenerateJWT(userID primitive.ObjectID, username, email, role, secret string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// validira JWT token i vraca podatke iz njega
// proverava potpis i da li je token jos uvek valjan
func ValidateJWT(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

