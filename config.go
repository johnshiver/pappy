package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBConn string `required:"true" split_words:"true" default:"postgres://jshiver@127.0.0.1:5432/jshiver"` // TODO: fix local connection string
}

func GetConfig() Config {
	var c Config
	envconfig.MustProcess("pappy", &c)
	return c
}

func GetDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", GetConfig().DBConn)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
