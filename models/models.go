package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

type Post struct {
	gorm.Model
	Title   string
	Content string
	File    sql.NullString
}
