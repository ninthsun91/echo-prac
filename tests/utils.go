package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"myapp/internal/db"
	"myapp/internal/lib/middlewares"
	"myapp/internal/routes"
	"net/http"
	"net/http/httptest"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitContext(method, target string, body interface{}) (echo.Context, *gorm.DB) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	db := db.ConnectDatabase()
	e.Use(middlewares.ContextDB(db))

	var req *http.Request
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			log.Fatalf("Error marshalling request body: %v", err)
		}
		req = httptest.NewRequest(method, target, bytes.NewBuffer(data))
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, db
}

func InitServer() *httptest.Server {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	db := db.ConnectDatabase()
	e.Use(middlewares.ContextDB(db))

	routes.Init(e.Group("api"), db)

	ts := httptest.NewServer(e)
	return ts
}
