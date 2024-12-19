// internal/services/user_service.go

// package services

// import (
// 	"github.com/Raviikumar001/e-com-api-go/internal/database"
// 	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
// )

// type UserService struct {
// 	userRepo repositories.UserRepository
// 	db       *database.DB
// }

// func NewUserService(userRepo repositories.UserRepository, db *database.DB) *UserService {
// 	return &UserService{userRepo: userRepo, db: db}
// }

// func (s *UserService) FindAll() ([]database.User, error) {
// 	return s.userRepo.FindAll()
// }

// func (s *UserService) FindByID(id uint) (database.User, error) {
// 	return s.userRepo.FindByID(id)
// }

// func (s *UserService) Create(user database.User) error {
// 	return s.userRepo.Create(user)
// }

// func (s *UserService) Update(user database.User) error {
// 	return s.userRepo.Update(user)
// }

// func (s *UserService) Delete(id uint) error {
// 	return s.userRepo.Delete(id)
// }

// internal/services/user_service.go

package services

import (
	"errors"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	userRepo    repositories.UserRepository
	authService *AuthService
}

func NewUserService(userRepo repositories.UserRepository, authService *AuthService) *UserService {
	return &UserService{userRepo: userRepo, authService: authService}
}

func (s *UserService) FindAll(r *fiber.Ctx) ([]database.User, error) {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return nil, err
	}
	if claims.Role != "admin" {
		return nil, errors.New("unauthorized")
	}
	return s.userRepo.FindAll()
}

func (s *UserService) FindByUsername(c *fiber.Ctx, username string) (database.User, error) {
	return s.userRepo.FindByUsername(username)
}

func (s *UserService) FindByID(r *fiber.Ctx, id uint) (database.User, error) {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return database.User{}, err
	}
	if claims.UserID != id {
		return database.User{}, errors.New("unauthorized")
	}
	return s.userRepo.FindByID(id)
}

func (s *UserService) Create(r *fiber.Ctx, user database.User) error {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return err
	}
	if claims.Role != "admin" {
		return errors.New("unauthorized")
	}
	return s.userRepo.Create(user)
}

func (s *UserService) Update(r *fiber.Ctx, user database.User) error {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return err
	}
	if claims.UserID != user.ID {
		return errors.New("unauthorized")
	}
	return s.userRepo.Update(user)
}

func (s *UserService) Delete(r *fiber.Ctx, id uint) error {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return err
	}
	if claims.UserID != id || claims.Role != "admin" {
		return errors.New("unauthorized")
	}
	return s.userRepo.Delete(id)
}
