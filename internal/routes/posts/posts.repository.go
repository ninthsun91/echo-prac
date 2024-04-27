package posts

import (
	"myapp/internal/db/models"

	"gorm.io/gorm"
)

type PostsRepository interface {
	Create(post models.Post) (models.Post, error)
	FindById(id uint) (models.Post, error)
	Update(id uint, data map[string]interface{}) (models.Post, error)
}

type postsRepository struct {
	db *gorm.DB
}

func NewPostsRepository(db *gorm.DB) PostsRepository {
	return &postsRepository{db}
}

func (r *postsRepository) Create(post models.Post) (models.Post, error) {
	result := r.db.Create(&post)
	if result.Error != nil {
		return models.Post{}, result.Error
	}
	return post, nil
}

func (r *postsRepository) FindById(id uint) (models.Post, error) {
	var post models.Post
	result := r.db.First(&post, id)
	if result.Error != nil {
		return models.Post{}, result.Error
	}
	return post, nil
}

func (r *postsRepository) Update(id uint, data map[string]interface{}) (models.Post, error) {
	var post models.Post
	findResult := r.db.Find(&post, id)
	if findResult.Error != nil {
		return models.Post{}, findResult.Error
	}
	updateResult := r.db.Model(&post).Updates(data)
	if updateResult.Error != nil {
		return models.Post{}, updateResult.Error
	}
	return post, nil
}
