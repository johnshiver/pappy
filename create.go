package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
)

func createDomain(db *gorm.DB, domainName, domainPassword, userPassword string, user *User) *Domain {
	key := generateEncryptionKey(user.Email, userPassword)
	enctypedDomainPw := encrypt(key, domainPassword)
	domainName = strings.ToLower(domainName)

	newDomain := Domain{
		FQDN:         domainName,
		PasswordHash: enctypedDomainPw,
		UserID:       user.ID,
	}
	err := db.Create(&newDomain).Error
	if err != nil {
		fmt.Println("there was an error creating db!")
		log.Panic(err)
	}
	fmt.Printf("Domain %s succesfully created!", newDomain.FQDN)
	return &newDomain
}
