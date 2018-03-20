package user

import (
	"fmt"
	"time"
)

const (
	UniqueConstraintUsername = "users_username_key"
	UniqueConstraintEmail    = "users_email_key"
)

type User struct {
	ID           uint   `gorm:"primary_key"`
	Email        string `gorm:"type:varchar(100);unique"`
	PasswordHash string
	CreatedAt    time.Time
}

type UsernameDuplicateError struct {
	Username string
}

func (e *UsernameDuplicateError) Error() string {
	return fmt.Sprintf("Username '%s' already exists", e.Username)
}
