// internal/services/auth_service.go

package services

import (
	"errors"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
	"github.com/Raviikumar001/e-com-api-go/internal/services/jwt_service"
)

type AuthService struct {
	userRepo   repositories.UserRepository
	jwtService jwt_service.JWTService
}

func NewAuthService(userRepo repositories.UserRepository, jwtService jwt_service.JWTService) *AuthService {
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

func (s *AuthService) ValidateToken(tokenString string) (jwt_service.JWTClaims, error) {
	return s.jwtService.ValidateToken(tokenString)
}
