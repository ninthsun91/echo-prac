package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"myapp/internal/db/models"
	"myapp/internal/routes/posts"

	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	ts, db := InitServer()
	defer ts.Close()

	url := fmt.Sprintf("%s/api/posts", ts.URL)

	cases := []struct {
		message string
		body    posts.PostCreateRequestBody
	}{
		{"empty request body", posts.PostCreateRequestBody{}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("[400]%s", tc.message), func(t *testing.T) {
			data, err := json.Marshal(tc.body)
			if err != nil {
				t.Fatalf("Error marshalling request body: %v", err)
			}

			resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
			if err != nil {
				t.Fatalf("Error sending POST request: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}

	t.Run("[201]post created", func(t *testing.T) {
		title := "Test Post"
		content := "This is a test post"
		test_file := "../public/test.txt"

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		err := writer.WriteField("title", title)
		assert.NoError(t, err)
		err = writer.WriteField("content", content)
		assert.NoError(t, err)

		file, err := os.Open(test_file)
		assert.NoError(t, err)
		defer file.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		assert.NoError(t, err)
		_, err = io.Copy(part, file)
		assert.NoError(t, err)
		writer.Close()

		resp, err := http.Post(url, writer.FormDataContentType(), body)
		if err != nil {
			t.Fatalf("Error sending POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		post := DecodeResBody[models.Post](t, resp)
		assert.NotNil(t, post)
		assert.Equal(t, title, post.Title)
		assert.Equal(t, content, post.Content)

		result := db.Unscoped().Delete(&post)
		assert.NoError(t, result.Error)
		fmt.Printf("Cleanup test post %d: %v", post.ID, result.Error)
	})
}

func TestFindOnePost(t *testing.T) {
	ts, _ := InitServer()
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

		post := DecodeResBody[models.Post](t, resp)
		assert.NotNil(t, post)
		assert.Equal(t, postId, post.ID)
	})
}
