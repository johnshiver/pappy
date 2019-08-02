package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/johnshiver/pappy/user"
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
	var lookup bool
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "add",
			Usage:       "add new object to database",
			Destination: &add,
		},
		cli.BoolFlag{
			Name:        "generate",
			Usage:       "generates object but doesnt persist",
			Destination: &generate,
		},
		cli.BoolFlag{
			Name:        "list",
			Usage:       "list objects in database",
			Destination: &list,
		},
		cli.BoolFlag{
			Name:        "lookup",
			Usage:       "lookup objects in database",
			Destination: &lookup,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "user",
			Usage: "Actions related to users",
			Action: func(c *cli.Context) error {
				if add == true {
					user.CreateUser(db)
				}

				return nil
			},
		},
		{
			Name:  "password",
			Usage: "takes no arguments",
			Action: func(c *cli.Context) error {
				if generate == true {
					pw := user.GeneratePassword()
					fmt.Println(pw)
				}
				return nil
			},
		},
		{
			Name:  "domain",
			Usage: "--create, --list, --lookup",
			Action: func(c *cli.Context) error {
				if add == true {
					user.CreateDomain(db)
				} else if list == true {
					user.ListDomains(db)
				} else if lookup == true {
					user.LookupDomain(db)
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
