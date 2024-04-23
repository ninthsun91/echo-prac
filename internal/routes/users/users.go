package users

import "github.com/labstack/echo/v4"

type UsersRouter struct{}

func (controller UsersRouter) Init(g *echo.Group) {
	g.POST("", controller.Signup)
}
