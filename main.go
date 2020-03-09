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

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
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
	db := config.GetDB()
	defer db.Close()

	env := initEnv(db)
	env.createTables()

	tgtObject := flag.String("object", "", "object to perform action on: user / password")
	action := flag.String("action", "", "action to perform on object: create / list / delete")

	flag.Parse()

	switch *tgtObject {
	case "user":
		switch *action {
		case "create":
			env.CreateUser()
		case "list":
			// this is a security hazard, I dont really care for my personal use
			env.ListUsers()
		}
	case "password":
		switch *action {
		case "create":
			env.LogIn()
			env.CreatePassword()
		case "list":
			env.LogIn()
			env.ListPasswords()
		case "delete":
			env.LogIn()
			env.DeletePassword()
		}
	default:
		switch *action {
		case "gen", "generate", "g":
			fmt.Println(generatePassword(-1))
		}
	}

}
