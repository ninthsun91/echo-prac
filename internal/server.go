package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"myapp/internal/db"
	"myapp/internal/db/models"
	"myapp/internal/validator"
)

type Server struct {
	e  *echo.Echo
	db *gorm.DB
}

func NewServer() *Server {
	e := echo.New()
	db := db.ConnectDatabase()

	return &Server{e, db}
}

func (s *Server) Start(addr string) {
	s.e.Validator = validator.SetCustomValidator()

	s.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	s.e.POST("/users", s.signupHandler)

	s.e.Logger.Fatal(s.e.Start(addr))
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

func (s *Server) signupHandler(c echo.Context) error {
	body := new(SignupRequestBody)
	if err := c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(body); err != nil {
		c.Logger().Errorf("Failed to validate request body: %v", err)
		return err
	}

	user := body.toUser()
	result := s.db.Create(&user)
	if result.Error != nil {
		c.Logger().Errorf("Failed to create user: %v", result.Error)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, user)
}
