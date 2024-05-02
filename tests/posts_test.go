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

	var postId uint = 1
	url := fmt.Sprintf("%s/api/posts/%d", ts.URL, postId)
	resp, err := http.Get(url)
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
}
