// internal/services/role_service.go

package services

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
)

type RoleService struct {
	roleRepo repositories.RoleRepository
	db       *database.DB
}

func NewRoleService(roleRepo repositories.RoleRepository, db *database.DB) *RoleService {
	return &RoleService{roleRepo: roleRepo, db: db}
}

func (s *RoleService) FindAll() ([]database.Role, error) {
	return s.roleRepo.FindAll()
}

func (s *RoleService) FindByID(id uint) (database.Role, error) {
	return s.roleRepo.FindByID(id)
}

func (s *RoleService) Create(role database.Role) error {
	return s.roleRepo.Create(role)
}

func (s *RoleService) Update(role database.Role) error {
	return s.roleRepo.Update(role)
}

func (s *RoleService) Delete(id uint) error {
	return s.roleRepo.Delete(id)
}
