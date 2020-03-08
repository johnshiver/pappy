package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"

	"github.com/johnshiver/pappy/config"
)

type runEnv struct {
	db            *sqlx.DB
	user          *User
	encryptionKey []byte

	mtx sync.Mutex
}

func initEnv() *runEnv {
	var env runEnv
	env.db = config.GetDB()
	return &env
}

func (env *runEnv) createTables() {
	env.CreateUserTable()
	env.CreatePasswordsTable()
}

func main() {
	env := initEnv()
	defer env.db.Close()

	env.createTables()

	if len(os.Args[1:]) != 1 {
		panic(fmt.Errorf("expected 1 cmd arg, received %d", len(os.Args[1:])))
	}

	cmdName := os.Args[1]
	switch cmdName {
	case "new_user":
		env.CreateUser()
	case "list":
		env.LogIn()
		env.ListPasswords()
	case "add":
		env.LogIn()
		env.CreatePassword()
	case "delete":
		env.LogIn()
	case "generate":
		generatePassword(-1)
	}

}
