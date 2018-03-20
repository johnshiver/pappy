package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/johnshiver/password_manager/user"
	"github.com/urfave/cli"
)

var (
	db_user     string
	db_password string
	db_name     string
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func getDbCreds() {
	db_user = os.Getenv("PW_MAN_DB_USER")
	db_password = os.Getenv("PW_MAN_DB_PW")
	db_name = os.Getenv("PW_MAN_DB_NAME")
}

func initDb() *gorm.DB {
	getDbCreds()
	db_string := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", db_user, db_name, db_password)
	db, err := gorm.Open("postgres", db_string)
	panicOnError(err)
	db.AutoMigrate(&user.User{}, &user.Domain{})
	return db
}

func main() {
	app := cli.NewApp()
	db := initDb()
	defer db.Close()

	var add bool
	var generate bool
	var list bool
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "add",
			Usage:       "add new object",
			Destination: &add,
		},
		cli.BoolFlag{
			Name:        "generate",
			Usage:       "add new object",
			Destination: &generate,
		},
		cli.BoolFlag{
			Name:        "list",
			Usage:       "list objects",
			Destination: &list,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "user",
			Aliases: []string{"c"},
			Usage:   "Actions related to users",
			Action: func(c *cli.Context) error {
				if add == true {
					new_user := user.CreateUser(db)
					fmt.Println(new_user)
				}

				return nil
			},
		},
		{
			Name:    "password",
			Aliases: []string{"a"},
			Usage:   "takes no arguments",
			Action: func(c *cli.Context) error {
				if generate == true {
					pw := user.GeneratePassword()
					fmt.Println(pw)
				}
				return nil
			},
		},
		{
			Name:    "domain",
			Aliases: []string{"a"},
			Usage:   "takes create and list",
			Action: func(c *cli.Context) error {
				if add == true {
					domain := user.CreateDomain(db)
					fmt.Println(domain)

				} else if list == true {
					user.ListDomains(db)
				}

				return nil
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
