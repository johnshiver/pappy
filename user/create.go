package user

import (
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

func Create(db *gorm.DB, user *User) (uint, error) {
	err := db.Create(user).Error
	if err != nil {
		if postgres.IsUniqueConstraintError(err, UniqueConstraintUsername) {
			return 0, &UsernameDuplicateError{Username: user.Email}
		}
		return 0, err
	}
	return user.ID, nil
}
