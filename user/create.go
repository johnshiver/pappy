package user

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sethvargo/go-diceware/diceware"
	"github.com/tamizhvendan/gomidway/postgres"
)

func createPassword() string {
	list, err := diceware.Generate(6)
	if err != nil {
		log.Fatal(err)
	}
	var upperList []string
	var randInt int
	for _, word := range list {
		word = strings.Title(word)
		randInt = rand.Intn(9)
		word += strconv.Itoa(randInt)
		upperList = append(upperList, strings.Title(word))
	}

	return strings.Join(upperList, "_")
}

func createUser(db *gorm.DB, user *User) (uint, error) {
	err := db.Create(user).Error
	if err != nil {
		if postgres.IsUniqueConstraintError(err, UniqueConstraintUsername) {
			return 0, &UsernameDuplicateError{Username: user.Email}
		}
		return 0, err
	}
	return user.ID, nil
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
