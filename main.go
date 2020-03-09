package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/johnshiver/pappy/config"
)

type runEnv struct {
	db        *sqlx.DB
	userInput io.Reader

	user          *User
	encryptionKey []byte

	mtx sync.Mutex
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func initEnv() *runEnv {
	var env runEnv
	env.db = config.GetDB()
	env.userInput = os.Stdin
	env.createTables()
	return &env
}

func (env *runEnv) createTables() {
	env.createUserTable()
	env.createPasswordsTable()
}

func main() {
	env := initEnv()
	defer env.db.Close()

	if len(os.Args[1:]) != 1 {
		panic(fmt.Errorf("expected 1 cmd arg, received %d", len(os.Args[1:])))
	}

	cmdName := os.Args[1]
	switch cmdName {
	case "new_user", "nu":
		env.CreateUser()
	case "list":
		env.LogIn()
		env.ListPasswords()
	case "add":
		env.LogIn()
		env.CreatePassword()
	case "delete", "del", "d":
		env.LogIn()
	case "generate", "gen", "g":
		fmt.Println(generatePassword(-1))
	}

}
