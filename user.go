package main

import (
	"database/sql"
	"time"
)

type User struct {
	ID           uint
	UserName     string
	PasswordHash string
	Passwords    []Password

	CreatedAt time.Time
}

func (env *runEnv) CreateUserTable() {
	createSQL := `
        CREATE TABLE if not exists users (
            id INTEGER PRIMARY KEY,
            user_name TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
		)
	`
	env.db.MustExec(createSQL)
}

func (env *runEnv) PersistUser(u User) {
	const insertSQL = `
           INSERT into users (user_name, password_hash)
           VALUES (LOWER($1), $2)
	`
	env.db.MustExec(insertSQL, u.UserName, u.PasswordHash)
}

func (env *runEnv) FindByUserName(userName string) (*User, error) {
	const userSQL = `
       SELECT ID, user_name, password_hash, created_at
       FROM users
       WHERE user_name=$1
       LIMIT 1
    `
	var user User
	err := env.db.Get(&user, userSQL, userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
