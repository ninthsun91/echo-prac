package users

import (
	"myapp/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

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

func (UsersRouter) Signup(c echo.Context) error {
	body := new(SignupRequestBody)
	if err := c.Bind(&body); err != nil {
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
