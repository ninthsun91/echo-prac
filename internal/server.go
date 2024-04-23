package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"myapp/internal/db"
	"myapp/internal/lib/middlewares"
	"myapp/internal/lib/validator"
	"myapp/internal/routes"
)

func Init(addr string) {
	e := echo.New()

	e.Validator = validator.SetCustomValidator()

	db := db.ConnectDatabase()
	e.Use(middlewares.ContextDB(db))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	routes.Routes(e.Group("api"))

	e.Logger.Fatal(e.Start(addr))
}
