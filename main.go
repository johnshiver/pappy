package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/johnshiver/pappy/config"
)

type runEnv struct {
	userInput io.Reader
	userSvc   UserService
	pwdSvc    PasswordService

	user          *User
	encryptionKey []byte

	mtx sync.Mutex
}


func initEnv(db *sqlx.DB) *runEnv {
	var env runEnv
	env.userInput = os.Stdin
	env.userSvc = NewUserDao(db)
	env.pwdSvc = NewPasswordDao(db)
	return &env
}

func (env *runEnv) createTables() {
	env.userSvc.CreateUserTable()
	env.pwdSvc.CreatePasswordsTable()
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	db := config.GetDB()
	defer db.Close()

	env := initEnv(db)
	env.createTables()

	tgtObject := flag.String("o", "", "object to perform action on: user / password")
	action := flag.String("a", "", "action to perform on object: create / list / delete")

	flag.Parse()

	switch *tgtObject {
	case "user", "u":
		switch *action {
		case "create", "c":
			env.CreateUser()
		case "list", "l":
			// this is a security hazard, I dont really care for my personal use
			env.ListUsers()
		}
	case "password", "p":
		switch *action {
		case "create", "c":
			env.LogIn()
			env.CreatePassword()
		case "list", "l":
			env.LogIn()
			env.ListPasswords()
		case "delete", "d":
			env.LogIn()
			env.DeletePassword()
		}
	default:
		switch *action {
		case "gen", "generate", "g":
			fmt.Println(generatePassword(-1))
		case "db", "dbloc":
			fmt.Println(config.GetDBFilePath())
		}
	}

}
