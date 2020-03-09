package main

import (
	"database/sql"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-diceware/diceware"
	log "github.com/sirupsen/logrus"
)

type Password struct {
	ID           uint
	Location     string
	PasswordHash string `db:"password_hash"`
	UserID       uint

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PasswordDao struct {
	db *sqlx.DB
}

func NewPasswordDao(db *sqlx.DB) *PasswordDao {
	return &PasswordDao{db: db}
}

type PasswordService interface {
	CreatePasswordsTable()
	CreatePassword(pw *Password)
	GetPasswords(userID uint) []*Password
	DeletePassword(userID uint, pwLoc string)
}

func (pd *PasswordDao) CreatePasswordsTable() {
	createSQL := `
        CREATE TABLE if not exists passwords (
            id INTEGER PRIMARY KEY,
            location TEXT NOT NULL,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                
            user_id INTEGER NOT NULL,
                                             
           FOREIGN KEY(user_id) REFERENCES users(id),
           UNIQUE(user_id, location)
		);
	`
	pd.db.MustExec(createSQL)
}

func (pd *PasswordDao) CreatePassword(pw *Password) {
	const insertSQL = `
           INSERT into passwords (location, password_hash, user_id)
           VALUES (LOWER($1), $2, $3)
	`
	pd.db.MustExec(insertSQL, pw.Location, pw.PasswordHash, pw.UserID)
}

func (pd *PasswordDao) GetPasswords(userID uint) []*Password {
	const pwSQL = `
         SELECT location, password_hash
         FROM passwords
         WHERE user_id=$1
    `
	var pws []*Password
	err := pd.db.Select(&pws, pwSQL, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}
	return pws
}

func (pd *PasswordDao) DeletePassword(userID uint, pwLoc string) {
	const deleteSQL = `
         DELETE FROM passwords WHERE user_id=$1 AND location=lower($2)
    `
	pd.db.MustExec(deleteSQL, userID, pwLoc)
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
