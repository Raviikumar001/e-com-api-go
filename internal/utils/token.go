// internal/utils/token.go
package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims struct for JWT
type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uint) (string, error) {
    // Get secret key from environment variable, or use a default for development
    secretKey := os.Getenv("JWT_SECRET_KEY")
    if secretKey == "" {
        secretKey = "your-secret-key-here" // For development only, in production always use environment variable
    }

    // Create claims with user ID and expiration time
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign the token with the secret key
    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// ValidateToken validates the JWT token and returns the user ID
func ValidateToken(tokenString string) (uint, error) {
    secretKey := os.Getenv("JWT_SECRET_KEY")
    if secretKey == "" {
        secretKey = "your-secret-key-here" // For development only
    }

    // Parse the token
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })

    if err != nil {
        return 0, err
    }

    // Check if token is valid
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims.UserID, nil
    }

    return 0, jwt.ErrSignatureInvalid
}