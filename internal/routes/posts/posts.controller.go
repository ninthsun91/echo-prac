package posts

import (
	"myapp/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PostsController struct {
	repo PostsRepository
}

func NewPostsController(repo PostsRepository) *PostsController {
	return &PostsController{repo}
}

func (posts *PostsController) Create(c echo.Context) error {
	var body PostCreateRequestBody
	if err := c.Bind(&body); err != nil {
		c.Logger().Errorf("Failed to bind request body: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(body); err != nil {
		c.Logger().Errorf("Failed to validate request body: %v", err)
		return err
	}

	post, err := posts.repo.Create(body.toPost(1))
	if err != nil {
		c.Logger().Errorf("Failed to create post: %v", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, post)
}

type PostCreateRequestBody struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (b PostCreateRequestBody) toPost(userId uint) models.Post {
	return models.Post{
		Title:   b.Title,
		Content: b.Content,
		UserID:  userId,
	}
}
