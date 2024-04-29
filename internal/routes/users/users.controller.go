package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"myapp/internal/db/models"
)

type UsersController struct {
	repo UsersRepository
}

func NewUsersController(repo UsersRepository) *UsersController {
	return &UsersController{repo}
}

func (uc *UsersController) Signup(c echo.Context) error {
	var body SignupRequestBody
	if err := c.Bind(&body); err != nil {
		c.Logger().Errorf("Failed to bind request body: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(body); err != nil {
		c.Logger().Errorf("Failed to validate request body: %v", err)
		return err
	}

	user, err := uc.repo.Create(body.toUser())
	if err != nil {
		c.Logger().Errorf("Failed to create user: %v", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, user)
}

func (uc *UsersController) FindUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Errorf("Invalid user ID: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	user, err := uc.repo.FindById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.String(http.StatusNotFound, "User not found")
		}

		c.Logger().Errorf("Failed to find user: %v", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, user)
}

func (uc *UsersController) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Errorf("Invalid user ID: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	var body UpdateUserRequestBody
	if err := c.Bind(&body); err != nil {
		c.Logger().Errorf("Failed to bind request body: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(body); err != nil {
		c.Logger().Errorf("Failed to validate request body: %v", err)
		return err
	}

	user, err := uc.repo.Update(uint(id), body.toMap())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.String(http.StatusNotFound, "User not found")
		}
		c.Logger().Errorf("Failed to update user: %v", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, user)
}

func (uc *UsersController) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Errorf("Invalid user ID: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	err = uc.repo.Delete(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.String(http.StatusNotFound, "User not found")
		}
		c.Logger().Errorf("Failed to delete user: %v", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.NoContent(http.StatusNoContent)
}

type SignupRequestBody struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

func (b *SignupRequestBody) toUser() models.User {
	return models.User{
		Name:     b.Name,
		Email:    b.Email,
		Password: b.Password,
	}
}

type UpdateUserRequestBody struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty" validate:"omitempty,min=6,max=20"`
}

func (b *UpdateUserRequestBody) toMap() map[string]interface{} {
	data := map[string]interface{}{}

	if b.Name != "" {
		data["name"] = b.Name
	}
	if b.Password != "" {
		data["password"] = b.Password
	}

	return data
}
