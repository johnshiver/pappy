package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
)

type runEnv struct {
	db            *sqlx.DB
	user          *User
	encryptionKey []byte

	mtx sync.Mutex
}

func initEnv() *runEnv {
	var env runEnv
	env.db = GetDB()
	return &env
}

func main() {
	env := initEnv()

	if len(os.Args[1:]) != 1 {
		panic(fmt.Errorf("expected 1 cmd arg, received %d", len(os.Args[1:])))
	}

	cmdName := os.Args[1]
	switch cmdName {
	case "list":
		env.LogIn()
		env.ListDomains()
	case "add":
		env.LogIn()
		env.CreatePassword()
	case "delete":
		env.LogIn()
	}

}
