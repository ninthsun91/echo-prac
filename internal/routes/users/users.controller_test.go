package users

import (
	"myapp/internal/db/models"
	"myapp/internal/routes/users/mocks"
	"myapp/internal/utils"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserSignup(t *testing.T) {
	t.Run("Invalid request body", func(t *testing.T) {
		body := strings.NewReader(`invalid json`)
		c, rec := utils.InitContext(http.MethodPost, "/users", body)

		mockRepo := new(mocks.MockUsersRepository)
		controller := NewUsersController(mockRepo)

		if assert.NoError(t, controller.Signup(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "Bad Request", rec.Body.String())
		}
	})

	t.Run("Request body validation failed", func(t *testing.T) {
		tests := []struct {
			name string
			body SignupRequestBody
		}{
			{
				name: "Empty name",
				body: SignupRequestBody{
					Email:    "test@meshed3d.com",
					Password: "qwe123",
				},
			},
			{
				name: "Empty email",
				body: SignupRequestBody{
					Name:     "Test User",
					Password: "qwe123",
				},
			},
			{
				name: "Invalid email format",
				body: SignupRequestBody{
					Name:     "Test User",
					Email:    "invalid-email",
					Password: "qwe123",
				},
			},
			{
				name: "Empty password",
				body: SignupRequestBody{
					Name:  "Test User",
					Email: "test@meshed3d.com",
				},
			},
			{
				name: "Short password",
				body: SignupRequestBody{
					Name:     "Test User",
					Email:    "test@meshed3d.com",
					Password: "qwe",
				},
			},
			{
				name: "Long password",
				body: SignupRequestBody{
					Name:     "Test User",
					Email:    "test@meshed3d.com",
					Password: "qwe12345678901234567890",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body := utils.EncodeReqBody(t, tt.body)
				c, rec := utils.InitContext(http.MethodPost, "/users", body)

				mockRepo := new(mocks.MockUsersRepository)
				controller := NewUsersController(mockRepo)

				if assert.NoError(t, controller.Signup(c)) {
					assert.Equal(t, http.StatusBadRequest, rec.Code)
				}
			})
		}
	})

	t.Run("Repository error - user creation failed", func(t *testing.T) {
		body := utils.EncodeReqBody(t, SignupRequestBody{
			Email:    "test@meshed3d.com",
			Name:     "Test User",
			Password: "qwe123",
		})
		c, rec := utils.InitContext(http.MethodPost, "/users", body)

		mockRepo := new(mocks.MockUsersRepository)
		mockRepo.On("Create", mock.Anything).Return(models.User{}, gorm.ErrDuplicatedKey).Once()
		controller := NewUsersController(mockRepo)

		if assert.NoError(t, controller.Signup(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("User created", func(t *testing.T) {
		body := utils.EncodeReqBody(t, SignupRequestBody{
			Email:    "test@meshed3d.com",
			Name:     "Test User",
			Password: "qwe123",
		})
		c, rec := utils.InitContext(http.MethodPost, "/users", body)

		mockRepo := new(mocks.MockUsersRepository)

		user := models.User{
			Model:    gorm.Model{ID: 1},
			Email:    "test@meshed3d.com",
			Name:     "Test User",
			Password: "qwe123",
		}
		mockRepo.On("Create", mock.Anything).Return(user, nil).Once()

		controller := NewUsersController(mockRepo)

		if assert.NoError(t, controller.Signup(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			data := utils.DecodeRecBody[models.User](t, rec)
			assert.Equal(t, user.ID, data.ID)
			assert.Equal(t, user.Email, data.Email)
			assert.Equal(t, user.Name, data.Name)
			assert.Equal(t, user.Password, data.Password)
		}
	})
}
