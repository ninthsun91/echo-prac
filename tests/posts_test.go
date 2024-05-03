package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"myapp/internal/db/models"

	"github.com/stretchr/testify/assert"
)

func TestFindOne(t *testing.T) {
	ts := InitServer()
	defer ts.Close()

	setUrl := func(id uint) string {
		return fmt.Sprintf("%s/api/posts/%d", ts.URL, id)
	}

	t.Run("400 - invalid post ID", func(t *testing.T) {
		resp, err := http.Get(setUrl(0))
		if err != nil {
			t.Fatalf("Error sending GET request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("404 - post not found", func(t *testing.T) {
		postId := uint(999)
		resp, err := http.Get(setUrl(postId))
		if err != nil {
			t.Fatalf("Error sending GET request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("200 - post found", func(t *testing.T) {
		postId := uint(1)
		resp, err := http.Get(setUrl(postId))
		if err != nil {
			t.Fatalf("Error sending GET request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var post models.Post
		err = json.NewDecoder(resp.Body).Decode(&post)
		assert.NoError(t, err)
		assert.NotNil(t, post)
		assert.Equal(t, postId, post.ID)
	})
}
