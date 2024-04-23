package routes

import (
	"myapp/internal/routes/users"

	"github.com/labstack/echo/v4"
)

func Routes(g *echo.Group) {
	users.UsersRouter{}.Init(g.Group("/users"))
}
