package main

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type User struct {
	ID           uint
	UserName     string `db:"user_name"`
	PasswordHash string `db:"password_hash"`

	CreatedAt time.Time `db:"created_at"`
}

type UserDao struct {
	db *sqlx.DB
}

func NewUserDao(db *sqlx.DB) *UserDao {
	return &UserDao{db: db}
}

type UserService interface {
	CreateUserTable()
	CreateUser(u User)
	GetUsers() []*User
	FindByUsername(userName string) (*User, error)
}

func (ud *UserDao) CreateUserTable() {
	createSQL := `
        CREATE TABLE if not exists users (
            id INTEGER PRIMARY KEY,
            user_name TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
		)
	`
	ud.db.MustExec(createSQL)
}

func (ud *UserDao) CreateUser(u User) {
	const insertSQL = `
           INSERT into users (user_name, password_hash)
           VALUES (LOWER($1), $2)
	`
	ud.db.MustExec(insertSQL, u.UserName, u.PasswordHash)
}

func (ud *UserDao) FindByUsername(userName string) (*User, error) {
	const userSQL = `
       SELECT ID, user_name, password_hash, created_at
       FROM users
       WHERE user_name=$1
       LIMIT 1
    `
	var user User
	err := ud.db.Get(&user, userSQL, userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ud *UserDao) GetUsers() []*User {
	const userSQL = `
         SELECT user_name
         FROM users
    `
	var users []*User
	err := ud.db.Select(&users, userSQL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal(err)
	}
	return users
}
