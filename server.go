package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"myapp/database"
)

func main() {
	e := echo.New()
	database.ConnectDatabase()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
