package posts

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func MapRoutes(g *echo.Group, db *gorm.DB) {
	repository := NewPostsRepository(db)
	controller := NewPostsController(repository)

	g.POST("", controller.Create)
	g.GET("/:id", controller.FindOne)
	g.PATCH("/:id", controller.Update)
	g.DELETE("/:id", controller.Delete)
}
