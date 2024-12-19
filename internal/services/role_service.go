// internal/services/role_service.go

// package services

// import (
// 	"github.com/Raviikumar001/e-com-api-go/internal/database"
// 	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
// )

// type RoleService struct {
// 	roleRepo repositories.RoleRepository
// 	db       *database.DB
// }

// func NewRoleService(roleRepo repositories.RoleRepository, db *database.DB) *RoleService {
// 	return &RoleService{roleRepo: roleRepo, db: db}
// }

// func (s *RoleService) FindAll() ([]database.Role, error) {
// 	return s.roleRepo.FindAll()
// }

// func (s *RoleService) FindByID(id uint) (database.Role, error) {
// 	return s.roleRepo.FindByID(id)
// }

// func (s *RoleService) Create(role database.Role) error {
// 	return s.roleRepo.Create(role)
// }

// func (s *RoleService) Update(role database.Role) error {
// 	return s.roleRepo.Update(role)
// }

// func (s *RoleService) Delete(id uint) error {
// 	return s.roleRepo.Delete(id)
// }

// internal/services/role_service.go

package services

import (
	"errors"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type RoleService struct {
	roleRepo    repositories.RoleRepository
	authService *AuthService
}

func NewRoleService(roleRepo repositories.RoleRepository, authService *AuthService) *RoleService {
	return &RoleService{roleRepo: roleRepo, authService: authService}
}

func (s *RoleService) FindAll(r *fiber.Ctx) ([]database.Role, error) {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return nil, err
	}
	if claims.Role != "admin" {
		return nil, errors.New("unauthorized")
	}
	return s.roleRepo.FindAll()
}

func (s *RoleService) FindByID(r *fiber.Ctx, id uint) (database.Role, error) {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return database.Role{}, err
	}
	if claims.Role != "admin" {
		return database.Role{}, errors.New("unauthorized")
	}
	return s.roleRepo.FindByID(id)
}

func (s *RoleService) Create(r *fiber.Ctx, role database.Role) error {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return err
	}
	if claims.Role != "admin" {
		return errors.New("unauthorized")
	}
	return s.roleRepo.Create(role)
}

func (s *RoleService) Update(r *fiber.Ctx, role database.Role) error {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return err
	}
	if claims.Role != "admin" {
		return errors.New("unauthorized")
	}
	return s.roleRepo.Update(role)
}

func (s *RoleService) Delete(r *fiber.Ctx, id uint) error {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return err
	}
	if claims.Role != "admin" {
		return errors.New("unauthorized")
	}
	return s.roleRepo.Delete(id)
}
