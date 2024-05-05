package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"myapp/internal/db/models"
	"myapp/internal/routes/users"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testuser models.User
var ts *httptest.Server

func TestMain(m *testing.M) {
	var db *gorm.DB
	ts, db = InitServer()
	code := m.Run()
	teardown(ts, db)
	os.Exit(code)
}

func teardown(ts *httptest.Server, db *gorm.DB) {
	result := db.Unscoped().Delete(&testuser)
	fmt.Printf("Cleanup test user %d: %v", testuser.ID, result.Error)
	ts.Close()
}

func TestSignUpE2e(t *testing.T) {
	url := fmt.Sprintf("%s/api/users", ts.URL)

	cases := []struct {
		message string
		form    users.SignupRequestBody
	}{
		{"empty request body", users.SignupRequestBody{}},
		{"missing password", users.SignupRequestBody{Name: "John", Email: "John@test.com"}},
		{"missing name", users.SignupRequestBody{Email: "John@test.com", Password: "1234"}},
		{"missing email", users.SignupRequestBody{Name: "John", Password: "1234"}},
		{"invalid email format", users.SignupRequestBody{Name: "John", Email: "John@test", Password: "1234"}},
		{"invalid email format", users.SignupRequestBody{Name: "John", Email: "John Doe", Password: "1234"}},
		{"password too short", users.SignupRequestBody{Name: "John", Email: "John@test.com", Password: "1234"}},
		{"password too long", users.SignupRequestBody{Name: "John", Email: "John@test.com", Password: "123456789012345678901"}},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("[400]%s", tc.message), func(t *testing.T) {
			data, err := json.Marshal(tc.form)
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

	t.Run("201 - user created", func(t *testing.T) {
		form := users.SignupRequestBody{
			Name:     "John",
			Email:    "John@test.com",
			Password: "qwe1234",
		}
		body := EncodeReqBody(t, form)

		resp, err := http.Post(url, "application/json", body)
		if err != nil {
			t.Fatalf("Error sending POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		testuser = DecodeResBody[models.User](t, resp)
		assert.NotNil(t, testuser)
		assert.Equal(t, form.Name, testuser.Name)
		assert.Equal(t, form.Email, testuser.Email)
	})
}

func TestFindUserE2e(t *testing.T) {
	setUrl := func(id any) string {
		return fmt.Sprintf("%s/api/users/%v", ts.URL, id)
	}

	cases := []struct {
		message string
		id      any
	}{
		{"id is 0", 0},
		{"id is negative", -1},
		{"mixture of characters", "qwe123"},
		{"missing id", nil},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("[400]%s", tc.message), func(t *testing.T) {
			url := setUrl(tc.id)

			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("Request Error GET %s : %v", url, err)
			}
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}

	t.Run("[404]user not found", func(t *testing.T) {
		userId := uint(999)
		url := setUrl(userId)

		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("Request Error: GET %s - %v", url, err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("[200]user found", func(t *testing.T) {
		url := setUrl(testuser.ID)

		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("Request Error: GET %s - %v", url, err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		user := DecodeResBody[models.User](t, resp)
		assert.NotNil(t, user)
		assert.Equal(t, testuser.ID, user.ID)
		assert.Equal(t, testuser.Name, user.Name)
		assert.Equal(t, testuser.Email, user.Email)
	})
}
