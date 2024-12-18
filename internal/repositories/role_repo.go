// internal/repositories/role_repo.go

package repositories

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *database.DB
}

func NewRoleRepository(db *database.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindAll() ([]database.Role, error) {
	var roles []database.Role
	result := r.db.Find(&roles)
	return roles, result.Error
}

func (r *RoleRepository) FindByID(id uint) (database.Role, error) {
	var role database.Role
	result := r.db.First(&role, id)
	return role, result.Error
}

func (r *RoleRepository) Create(role database.Role) error {
	result := r.db.Create(&role)
	return result.Error
}

func (r *RoleRepository) Update(role database.Role) error {
	result := r.db.Save(&role)
	return result.Error
}

func (r *RoleRepository) Delete(id uint) error {
	result := r.db.Delete(&database.Role{Model: gorm.Model{ID: id}})
	return result.Error
}
