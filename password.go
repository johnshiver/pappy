package main

import (
	"context"
	"database/sql"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-diceware/diceware"
	log "github.com/sirupsen/logrus"
)

type PasswordService interface {
	CreatePasswordsTable(ctx context.Context, db *sqlx.DB)
	CreatePassword(ctx context.Context, db *sqlx.DB, pw *Password)
	GetPasswords(ctx context.Context, db *sqlx.DB, userID uint) []*Password
	DeletePassword(ctx context.Context, db *sqlx.DB, userID uint, pwLoc string)
}

type Password struct {
	ID           uint
	Location     string
	PasswordHash string `db:"password_hash"`
	UserID       uint

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PasswordDao struct {}

func NewPasswordDao() PasswordDao {
	return PasswordDao{}
}

const createPasswordsTableSQL = `
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

func (pd PasswordDao) CreatePasswordsTable(ctx context.Context, db *sqlx.DB) {
	db.MustExecContext(ctx, createPasswordsTableSQL)
}

const createPasswordSQL = `
INSERT into passwords (location, password_hash, user_id)
VALUES (LOWER($1), $2, $3)
`
func (pd PasswordDao) CreatePassword(ctx context.Context, db *sqlx.DB, pw *Password) {
	db.MustExecContext(ctx, createPasswordSQL, pw.Location, pw.PasswordHash, pw.UserID)
}

func (pd PasswordDao) GetPasswords(ctx context.Context, db *sqlx.DB, userID uint) []*Password {
	const pwSQL = `
         SELECT location, password_hash
         FROM passwords
         WHERE user_id=$1
    `
	var pws []*Password
	err := db.SelectContext(ctx, &pws, pwSQL, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}
	return pws
}

const deletePasswordSQL = `DELETE FROM passwords WHERE user_id=$1 AND location=lower($2)`

func (pd PasswordDao) DeletePassword(ctx context.Context, db *sqlx.DB, userID uint, pwLoc string) {
	db.MustExecContext(ctx, deletePasswordSQL, userID, pwLoc)
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

