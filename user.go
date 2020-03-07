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

func (env *runEnv) PersistUser(u *User) {

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
