package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/sethvargo/go-diceware/diceware"
)

type Password struct {
	ID           uint
	Location     string
	PasswordHash string
	CreatedAt    time.Time
	UserID       uint
}

func generatePassword(passwordLength int) string {
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

	newPassword := strings.Join(upperList, "_")
	if passwordLength < 0 {
		return newPassword
	}
	return newPassword[:passwordLength]
}

func (env *runEnv) createDomain(domainName, domainPassword, userPassword string, user *User) *Domain {

	// put SQL here

	key := generateEncryptionKey(user.ID, userPassword)
	enctypedDomainPw := encrypt(key, domainPassword)
	domainName = strings.ToLower(domainName)

	newDomain := Domain{
		FQDN:         domainName,
		PasswordHash: enctypedDomainPw,
		UserID:       user.ID,
	}
	if err != nil {
		fmt.Println("there was an error creating db!")
		log.Panic(err)
	}
	fmt.Printf("Domain %s succesfully created!", newDomain.FQDN)
	return &newDomain
}
