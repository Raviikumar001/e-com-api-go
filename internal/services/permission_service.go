// internal/services/permission_service.go

// package services

// import (
// 	"github.com/Raviikumar001/e-com-api-go/internal/database"
// 	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
// )

// type PermissionService struct {
// 	permissionRepo repositories.PermissionRepository
// 	db             *database.DB
// }

// func NewPermissionService(permissionRepo repositories.PermissionRepository, db *database.DB) *PermissionService {
// 	return &PermissionService{permissionRepo: permissionRepo, db: db}
// }

// func (s *PermissionService) FindAll() ([]database.Permission, error) {
// 	return s.permissionRepo.FindAll()
// }

// func (s *PermissionService) FindByID(id uint) (database.Permission, error) {
// 	return s.permissionRepo.FindByID(id)
// }

// func (s *PermissionService) Create(permission database.Permission) error {
// 	return s.permissionRepo.Create(permission)
// }

// func (s *PermissionService) Update(permission database.Permission) error {
// 	return s.permissionRepo.Update(permission)
// }

// func (s *PermissionService) Delete(id uint) error {
// 	return s.permissionRepo.Delete(id)
// }

// internal/services/permission_service.go

package services

import (
	"errors"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type PermissionService struct {
	permissionRepo repositories.PermissionRepository
	authService    *AuthService
}

func NewPermissionService(permissionRepo repositories.PermissionRepository, authService *AuthService) *PermissionService {
	return &PermissionService{permissionRepo: permissionRepo, authService: authService}
}

func (s *PermissionService) FindAll(r *fiber.Ctx) ([]database.Permission, error) {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return nil, err
	}
	if claims.Role != "admin" {
		return nil, errors.New("unauthorized")
	}
	return s.permissionRepo.FindAll()
}

func (s *PermissionService) FindByID(r *fiber.Ctx, id uint) (database.Permission, error) {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return database.Permission{}, err
	}
	if claims.Role != "admin" {
		return database.Permission{}, errors.New("unauthorized")
	}
	return s.permissionRepo.FindByID(id)
}

func (s *PermissionService) Create(r *fiber.Ctx, permission database.Permission) error {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return err
	}
	if claims.Role != "admin" {
		return errors.New("unauthorized")
	}
	return s.permissionRepo.Create(permission)
}

func (s *PermissionService) Update(r *fiber.Ctx, permission database.Permission) error {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return err
	}
	if claims.Role != "admin" {
		return errors.New("unauthorized")
	}
	return s.permissionRepo.Update(permission)
}

func (s *PermissionService) Delete(r *fiber.Ctx, id uint) error {
	claims, err := s.authService.ValidateToken(r)
	if err != nil {
		return err
	}
	if claims.Role != "admin" {
		return errors.New("unauthorized")
	}
	return s.permissionRepo.Delete(id)
}
