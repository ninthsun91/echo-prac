package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"myapp/database"
	"myapp/models"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		echo.New().AcquireContext().Logger().Errorf("Failed to validate request body: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	database.ConnectDatabase()

	e.Validator = &CustomValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/users", signupHandler)

	e.Logger.Fatal(e.Start(":1323"))
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

func signupHandler(c echo.Context) error {
	body := new(SignupRequestBody)
	if err := c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(body); err != nil {
		c.Logger().Errorf("Failed to validate request body: %v", err)
		return err
	}

	user := body.toUser()
	result := database.DB.Create(&user)
	if result.Error != nil {
		c.Logger().Errorf("Failed to create user: %v", result.Error)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, user)
}
