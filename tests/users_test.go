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

var userE2eEnv *UserE2e

type UserE2e struct {
	testuser models.User
	ts       *httptest.Server
	db       *gorm.DB
}

func TestMain(m *testing.M) {
	userE2eEnv = &UserE2e{}
	userE2eEnv.Setup()

	code := m.Run()

	userE2eEnv.teardown()
	os.Exit(code)
}

func (env *UserE2e) Setup() {
	env.ts, env.db = InitServer()
}

func (env *UserE2e) teardown() {
	result := env.db.Unscoped().Delete(&env.testuser)
	fmt.Printf("Cleanup test user %d: %v", env.testuser.ID, result.Error)
	env.ts.Close()
}

func TestSignup(t *testing.T) {
	testSignUp(t, userE2eEnv)
}
func testSignUp(t *testing.T, env *UserE2e) {
	url := fmt.Sprintf("%s/api/users", env.ts.URL)

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

		env.testuser = DecodeResBody[models.User](t, resp)
		assert.NotNil(t, env.testuser)
		assert.Equal(t, form.Name, env.testuser.Name)
		assert.Equal(t, form.Email, env.testuser.Email)
	})
}

func TestFindUser(t *testing.T) {
	testFindUser(t, userE2eEnv)
}
func testFindUser(t *testing.T, env *UserE2e) {
	setUrl := func(id any) string {
		return fmt.Sprintf("%s/api/users/%v", env.ts.URL, id)
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
		url := setUrl(env.testuser.ID)

		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("Request Error: GET %s - %v", url, err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		user := DecodeResBody[models.User](t, resp)
		assert.NotNil(t, user)
		assert.Equal(t, env.testuser.ID, user.ID)
		assert.Equal(t, env.testuser.Name, user.Name)
		assert.Equal(t, env.testuser.Email, user.Email)
	})
}

func TestUpdateUser(t *testing.T) {
	testUpdateUser(t, userE2eEnv)
}
func testUpdateUser(t *testing.T, env *UserE2e) {
	setUrl := func(id any) string {
		return fmt.Sprintf("%s/api/users/%v", env.ts.URL, id)
	}

	idCases := []struct {
		message string
		id      any
	}{
		{"id is 0", 0},
		{"id is negative", -1},
		{"mixture of characters", "qwe123"},
		{"missing id", nil},
	}
	for _, tc := range idCases {
		t.Run(fmt.Sprintf("[400]%s", tc.message), func(t *testing.T) {
			url := setUrl(tc.id)
			body := EncodeReqBody(t, users.UpdateUserRequestBody{})

			req, err := http.NewRequest(http.MethodPatch, url, body)
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatalf("Error creating PATCH request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Error sending PATCH request: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}

	bodyCases := []struct {
		message string
		form    users.UpdateUserRequestBody
	}{
		{"password too short", users.UpdateUserRequestBody{Name: "John", Password: "1234"}},
		{"password too long", users.UpdateUserRequestBody{Name: "John", Password: "123456789012345678901"}},
	}
	for _, tc := range bodyCases {
		t.Run(fmt.Sprintf("[400]%s", tc.message), func(t *testing.T) {
			url := setUrl(env.testuser.ID)
			body := EncodeReqBody(t, tc.form)

			req, err := http.NewRequest(http.MethodPatch, url, body)
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatalf("Error creating PATCH request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Error sending PATCH request: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}

	t.Run("[404]user not found", func(t *testing.T) {
		url := setUrl(uint(999))
		body := EncodeReqBody(t, users.UpdateUserRequestBody{Name: "Updated", Password: "qwe1234"})

		req, err := http.NewRequest(http.MethodPatch, url, body)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("Error creating PATCH request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Error sending PATCH request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("[200]user updated", func(t *testing.T) {
		url := setUrl(env.testuser.ID)
		form := users.UpdateUserRequestBody{Name: "Updated"}
		body := EncodeReqBody(t, form)

		req, err := http.NewRequest(http.MethodPatch, url, body)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("Error creating PATCH request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Error sending PATCH request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		user := DecodeResBody[models.User](t, resp)
		assert.Equal(t, env.testuser.ID, user.ID)
		assert.Equal(t, form.Name, user.Name)
	})
}
