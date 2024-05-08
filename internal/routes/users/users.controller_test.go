package users

import (
	"myapp/internal/db/models"
	"myapp/internal/routes/users/mocks"
	"myapp/internal/utils"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestSignup(t *testing.T) {
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
}
