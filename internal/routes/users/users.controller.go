package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"myapp/internal/db/models"
)

func (UsersRouter) Signup(c echo.Context) error {
	var body SignupRequestBody
	if err := c.Bind(&body); err != nil {
		c.Logger().Errorf("Failed to bind request body: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(body); err != nil {
		c.Logger().Errorf("Failed to validate request body: %v", err)
		return err
	}

	db := c.Get("db").(*gorm.DB)
	user := body.toUser()
	result := db.Create(&user)
	if result.Error != nil {
		c.Logger().Errorf("Failed to create user: %v", result.Error)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, user)
}

func (UsersRouter) FindUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Errorf("Invalid user ID: %v", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	db := c.Get("db").(*gorm.DB)
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.String(http.StatusNotFound, "User not found")
		}

		c.Logger().Errorf("Failed to find user: %v", result.Error)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, user)
}

func (UsersRouter) UpdateUser(c echo.Context) error {
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

	db := c.Get("db").(*gorm.DB)
	var user models.User

	findResult := db.First(&user, id)
	if findResult.Error != nil {
		if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
			return c.String(http.StatusNotFound, "User not found")
		}
		c.Logger().Errorf("Failed to find user: %v", findResult.Error)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	updateResult := db.Model(&user).Updates(body.toMap())
	if updateResult.Error != nil {
		c.Logger().Errorf("Failed to update user: %v", updateResult.Error)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, user)
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
