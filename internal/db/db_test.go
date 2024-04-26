package db

import (
	"os"
	"testing"

	"myapp/internal/db/models"
)

func TestConnectDatabase(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "qwe123123")
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_PORT", "5432")
	db := ConnectDatabase()

	if db == nil {
		t.Fatalf("DB is nil")
	}

	if !db.Migrator().HasTable(&models.User{}) {
		t.Fatalf("users table does not exist")
	}

	if !db.Migrator().HasTable(&models.Post{}) {
		t.Fatalf("posts table does not exist")
	}
}
