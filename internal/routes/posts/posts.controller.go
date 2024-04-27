package posts

import (
	"errors"
	"net/http"
	"strconv"

	"myapp/internal/db/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

func (posts *PostsController) FindOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if (err != nil) || (id < 1) {
		c.Logger().Errorf("Invalid post ID: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	post, err := posts.repo.FindById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.String(http.StatusNotFound, "Post not found")
		}
		c.Logger().Errorf("Failed to find post: %v", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, post)
}

func (posts *PostsController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if (err != nil) || (id < 1) {
		c.Logger().Errorf("Invalid post ID: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	var body PostUpdateRequestBody
	if err := c.Bind(&body); err != nil {
		c.Logger().Errorf("Failed to bind request body: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(body); err != nil {
		c.Logger().Errorf("Failed to validate request body: %v", err)
		return err
	}

	post, err := posts.repo.Update(uint(id), body.toPost())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.String(http.StatusNotFound, "Post not found")
		}
		c.Logger().Errorf("Failed to update post: %v", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, post)
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

type PostUpdateRequestBody struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

func (b PostUpdateRequestBody) toPost() map[string]interface{} {
	data := make(map[string]interface{})
	if b.Title != "" {
		data["title"] = b.Title
	}
	if b.Content != "" {
		data["content"] = b.Content
	}
	return data
}
