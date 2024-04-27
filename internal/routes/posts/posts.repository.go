package posts

import (
	"myapp/internal/db/models"

	"gorm.io/gorm"
)

type PostsRepository interface {
	Create(post models.Post) (models.Post, error)
}

type postsRepository struct {
	db *gorm.DB
}

func NewPostsRepository(db *gorm.DB) PostsRepository {
	return &postsRepository{db}
}

func (r *postsRepository) Create(post models.Post) (models.Post, error) {
	result := r.db.Create(&post)
	return post, result.Error
}
