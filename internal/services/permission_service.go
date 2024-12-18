// internal/services/permission_service.go

package services

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
)

type PermissionService struct {
	permissionRepo repositories.PermissionRepository
	db             *database.DB
}

func NewPermissionService(permissionRepo repositories.PermissionRepository, db *database.DB) *PermissionService {
	return &PermissionService{permissionRepo: permissionRepo, db: db}
}

func (s *PermissionService) FindAll() ([]database.Permission, error) {
	return s.permissionRepo.FindAll()
}

func (s *PermissionService) FindByID(id uint) (database.Permission, error) {
	return s.permissionRepo.FindByID(id)
}

func (s *PermissionService) Create(permission database.Permission) error {
	return s.permissionRepo.Create(permission)
}

func (s *PermissionService) Update(permission database.Permission) error {
	return s.permissionRepo.Update(permission)
}

func (s *PermissionService) Delete(id uint) error {
	return s.permissionRepo.Delete(id)
}
