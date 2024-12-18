// internal/repositories/user_repo.go

package repositories

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll() ([]database.User, error) {
	var users []database.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *UserRepository) FindByID(id uint) (database.User, error) {
	var user database.User
	result := r.db.First(&user, id)
	return user, result.Error
}

func (r *UserRepository) Create(user database.User) error {
	result := r.db.Create(&user)
	return result.Error
}

func (r *UserRepository) Update(user database.User) error {
	result := r.db.Save(&user)
	return result.Error
}

func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&database.User{Model: gorm.Model{ID: id}})
	return result.Error
}
