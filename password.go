package main

import (
	"database/sql"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/sethvargo/go-diceware/diceware"
	log "github.com/sirupsen/logrus"
)

type Password struct {
	ID           uint
	Location     string
	PasswordHash string
	UserID       uint

	CreatedAt time.Time
}

func (env *runEnv) createPasswordsTable() {
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
	if passwordLength <= 0 {
		return newPassword
	}
	return newPassword[:passwordLength]
}

func (env *runEnv) createPassword(pw *Password) {
	const insertSQL = `
           INSERT into passwords (location, password_hash, user_id)
           VALUES (LOWER($1), $2, $3)
	`
	env.db.MustExec(insertSQL, pw.Location, pw.PasswordHash, pw.UserID)
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
		log.Fatal(err)
	}
	return pws
}
