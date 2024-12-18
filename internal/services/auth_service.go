// internal/services/auth_service.go

package services

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
)

type AuthService struct {
	userRepo repositories.UserRepository
	db       *database.DB
}

func NewAuthService(userRepo repositories.UserRepository, db *database.DB) *AuthService {
	return &AuthService{userRepo: userRepo, db: db}
}

func (s *AuthService) Authenticate(username string, password string) (database.User, error) {
	// TO DO: implement authentication logic
	return database.User{}, nil
}
