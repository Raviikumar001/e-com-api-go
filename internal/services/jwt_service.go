// internal/services/jwt_service.go

package services

import (
	"errors"
	"time"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	secretKey string
	issuer    string
}

func NewJWTService(secretKey string, issuer string) *JWTService {
	return &JWTService{secretKey: secretKey, issuer: issuer}
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func (s *JWTService) GenerateToken(user database.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Role:   user.Role, // Assuming the user struct has a Role field
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer:    s.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	return claims, nil
}
