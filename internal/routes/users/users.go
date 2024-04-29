package users

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func MapRoutes(g *echo.Group, db *gorm.DB) {
	repository := NewUsersRepository(db)
	controller := NewUsersController(repository)

	g.POST("", controller.Signup)
	g.GET("/:id", controller.FindUser)
	g.PATCH("/:id", controller.UpdateUser)
	g.DELETE("/:id", controller.DeleteUser)
}
