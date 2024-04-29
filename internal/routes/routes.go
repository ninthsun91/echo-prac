package routes

import (
	"myapp/internal/routes/posts"
	"myapp/internal/routes/users"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Init(g *echo.Group, db *gorm.DB) {
	users.MapRoutes(g.Group("/users"), db)
	posts.MapRoutes(g.Group("/posts"), db)
}
