package db

import (
	"os"
	"testing"

	"myapp/internal/db/models"

	"github.com/stretchr/testify/assert"
)

func TestConnectDatabase(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "qwe123123")
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_PORT", "5432")

	db := ConnectDatabase()

	assert.NotNil(t, db, "db connection failed")

	hasUserTable := db.Migrator().HasTable(&models.User{})
	assert.True(t, hasUserTable, "missing user table")

	hasPostTable := db.Migrator().HasTable(&models.Post{})
	assert.True(t, hasPostTable, "missing post table")
}
