// internal/services/auth_service.go

package services

import (
	"errors"
	"net/http"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
)

type AuthService struct {
	userRepo   repositories.UserRepository
	jwtService *JWTService
}

func NewAuthService(userRepo repositories.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{userRepo: userRepo, jwtService: jwtService}
}

func (s *AuthService) Authenticate(username string, password string) (string, error) {
	var user database.User
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	if !user.ComparePassword(password) {
		return "", errors.New("invalid password")
	}
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) ValidateToken(r *http.Request) (*JWTClaims, error) {
	tokenString := r.Header.Get("Authorization")
	return s.jwtService.ValidateToken(tokenString)
}
