package users

import (
	"myapp/internal/db/models"

	"gorm.io/gorm"
)

type UsersRepository interface {
	Create(user models.User) (models.User, error)
	FindById(id uint) (models.User, error)
	Update(id uint, data map[string]interface{}) (models.User, error)
	Delete(id uint) error
}

type usersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &usersRepository{db}
}

func (r *usersRepository) Create(user models.User) (models.User, error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *usersRepository) FindById(id uint) (models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *usersRepository) Update(id uint, data map[string]interface{}) (models.User, error) {
	var user models.User
	findResult := r.db.First(&user, id)
	if findResult.Error != nil {
		return models.User{}, findResult.Error
	}
	updateResult := r.db.Model(&user).Updates(data)
	if updateResult.Error != nil {
		return models.User{}, updateResult.Error
	}
	return user, nil
}

func (r *usersRepository) Delete(id uint) error {
	var user models.User
	findResult := r.db.First(&user, id)
	if findResult.Error != nil {
		return findResult.Error
	}
	deleteResult := r.db.Delete(&user)
	return deleteResult.Error
}
