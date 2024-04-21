package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Posts    []Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Post struct {
	gorm.Model
	UserID  uint
	Title   string
	Content string
	File    sql.NullString
}
