package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myapp/internal/db/models"
	"myapp/internal/routes/users"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignUpE2e(t *testing.T) {
	ts := InitServer()
	defer ts.Close()

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
		data, err := json.Marshal(form)
		if err != nil {
			t.Fatalf("Error marshalling request body: %v", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("Error sending POST request: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var user models.User
		err = json.NewDecoder(resp.Body).Decode(&user)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, form.Name, user.Name)
		assert.Equal(t, form.Email, user.Email)
	})
}
