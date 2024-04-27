package routes

import (
	"myapp/internal/routes/posts"
	"myapp/internal/routes/users"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Routes(g *echo.Group, db *gorm.DB) {
	users.UsersRouter{}.Init(g.Group("/users"))
	posts.PostsRouter(g.Group("/posts"), db)
}
