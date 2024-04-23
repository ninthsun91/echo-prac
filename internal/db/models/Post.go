package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID  uint
	Title   string
	Content string
	File    sql.NullString
}
