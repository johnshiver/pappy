package main

import (
	"database/sql"
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

func (env *runEnv) CreatePasswordsTable() {
	createSQL := `
        CREATE TABLE if not exists passwords (
            id INTEGER PRIMARY KEY,
            location TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                
            user_id INTEGER NOT NULL,
                                             
           FOREIGN KEY(user_id) REFERENCES users(id)
		);
	`
	env.db.MustExec(createSQL)
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

func (env *runEnv) createPassword() error {

	// put SQL here

	//key := generateEncryptionKey(user.ID, userPassword)
	//enctypedDomainPw := encrypt(key, domainPassword)
	//domainName = strings.ToLower(domainName)
	//
	//newDomain := Domain{
	//	FQDN:         domainName,
	//	PasswordHash: enctypedDomainPw,
	//	UserID:       user.ID,
	//}
	//if err != nil {
	//	fmt.Println("there was an error creating db!")
	//	log.Panic(err)
	//}
	//fmt.Printf("Domain %s succesfully created!", newDomain.FQDN)
	return nil
}

func (env *runEnv) GetPasswords() []Password {
	const pwSQL = `
         SELECT *
         FROM passwords
         WHERE user_id=$1
    `
	var pws []Password
	err := env.db.Select(pws, pwSQL, env.user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	return pws
}
