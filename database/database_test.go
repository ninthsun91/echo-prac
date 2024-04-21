package database

import (
	"os"
	"testing"

	"myapp/models"
)

func TestConnectDatabase(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "qwe123123")
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_PORT", "5432")
	ConnectDatabase()

	if DB == nil {
		t.Fatalf("DB is nil")
	}

	if !DB.Migrator().HasTable(&models.User{}) {
		t.Fatalf("users table does not exist")
	}

	if !DB.Migrator().HasTable(&models.Post{}) {
		t.Fatalf("posts table does not exist")
	}
}
