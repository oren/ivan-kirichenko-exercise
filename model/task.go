package model

import "time"

// Task defines some todo-task to keep in our database
type Task struct {
	Id          int64 `gorm:"primary_key",sql:"AUTO_INCREMENT"`
	Title       string
	Description string
	Priority    int
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	CompletedAt *time.Time
	IsDeleted   bool
	IsCompleted bool
}
