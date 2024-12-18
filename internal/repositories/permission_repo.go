// internal/repositories/permission_repo.go

package repositories

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *database.DB
}

func NewPermissionRepository(db *database.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) FindAll() ([]database.Permission, error) {
	var permissions []database.Permission
	result := r.db.Find(&permissions)
	return permissions, result.Error
}

func (r *PermissionRepository) FindByID(id uint) (database.Permission, error) {
	var permission database.Permission
	result := r.db.First(&permission, id)
	return permission, result.Error
}

func (r *PermissionRepository) Create(permission database.Permission) error {
	result := r.db.Create(&permission)
	return result.Error
}

func (r *PermissionRepository) Update(permission database.Permission) error {
	result := r.db.Save(&permission)
	return result.Error
}

func (r *PermissionRepository) Delete(id uint) error {
	result := r.db.Delete(&database.Permission{Model: gorm.Model{ID: id}})
	return result.Error
}
