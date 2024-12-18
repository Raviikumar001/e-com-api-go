// internal/services/user_service.go

package services

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
	db       *database.DB
}

func NewUserService(userRepo repositories.UserRepository, db *database.DB) *UserService {
	return &UserService{userRepo: userRepo, db: db}
}

func (s *UserService) FindAll() ([]database.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) FindByID(id uint) (database.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) Create(user database.User) error {
	return s.userRepo.Create(user)
}

func (s *UserService) Update(user database.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}
