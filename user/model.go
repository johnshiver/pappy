package user

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
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
	Domains      []Domain
}

type Domain struct {
	ID           uint `gorm:"primary_key"`
	FQDN         string
	PasswordHash string
	CreatedAt    time.Time
	UserID       uint
}

func createDomain(db *gorm.DB, domain_name, domain_password, user_password string, user *User) *Domain {
	key := generateEncryptionKey(user.Email, user_password)
	enctyped_domain_pw := encrypt(key, domain_password)
	domain_name = strings.ToLower(domain_name)

	new_domain := Domain{
		FQDN:         domain_name,
		PasswordHash: enctyped_domain_pw,
		UserID:       user.ID,
	}
	err := db.Create(&new_domain).Error
	if err != nil {
		fmt.Println("there was an error creating db!")
		log.Panic(err)
	}
	fmt.Printf("Domain %s succesfully created!", new_domain.FQDN)
	return &new_domain
}

type UsernameDuplicateError struct {
	Username string
}

func (e *UsernameDuplicateError) Error() string {
	return fmt.Sprintf("Username '%s' already exists", e.Username)
}
