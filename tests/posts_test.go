package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"myapp/internal/db"
	"myapp/internal/db/models"
	"myapp/internal/lib/middlewares"
	"myapp/internal/routes/posts"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFindOne(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	e := echo.New()
	db := db.ConnectDatabase()
	e.Use(middlewares.ContextDB(db))

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var postId uint = 1
	c.SetPath("/posts/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatUint(uint64(postId), 10))

	repo := posts.NewPostsRepository(db)
	controller := posts.NewPostsController(repo)

	if assert.NoError(t, controller.FindOne(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var post models.Post
		err := json.NewDecoder(rec.Body).Decode(&post)
		assert.NoError(t, err)
		assert.NotNil(t, post)
		assert.Equal(t, postId, post.ID)
	}
}
