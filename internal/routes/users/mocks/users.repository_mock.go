package mocks

import (
	"myapp/internal/db/models"

	"github.com/stretchr/testify/mock"
)

type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) Create(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUsersRepository) FindById(id uint) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUsersRepository) Update(id uint, data map[string]interface{}) (models.User, error) {
	args := m.Called(id, data)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUsersRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
