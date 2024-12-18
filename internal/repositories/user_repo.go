// internal/repositories/user_repo.go

package repositories

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]database.User, error)
	FindByID(id uint) (database.User, error)
	FindByUsername(username string) (database.User, error)
	Create(user database.User) error
	Update(user database.User) error
	Delete(id uint) error
}

type userRepo struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindAll() ([]database.User, error) {
	var users []database.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *userRepo) FindByID(id uint) (database.User, error) {
	var user database.User
	result := r.db.First(&user, id)
	return user, result.Error
}

func (r *userRepo) FindByUsername(username string) (database.User, error) {
	var user database.User
	result := r.db.Where("username = ?", username).First(&user)
	return user, result.Error
}

func (r *userRepo) Create(user database.User) error {
	result := r.db.Create(&user)
	return result.Error
}

func (r *userRepo) Update(user database.User) error {
	result := r.db.Save(&user)
	return result.Error
}

func (r *userRepo) Delete(id uint) error {
	result := r.db.Delete(&database.User{Model: gorm.Model{ID: id}})
	return result.Error
}
