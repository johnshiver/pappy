package config

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var c *Config

type Config struct {
	DBDirName  string `required:"true" split_words:"true" default:".pappy"`
	DBFileName string `required:"true" split_words:"true" default:"data.db"`
}

func init() {
	c = &Config{}
	envconfig.MustProcess("PAPPY", c)
}

func GetDB() *sqlx.DB {
	checkOrCreatDBFiles()
	db, err := sqlx.Open("sqlite3", GetDBFilePath())
	if err != nil {
		log.Panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Panic(err)
	}
	return db
}

func GetDBFilePath() string {
	return getDBDirPath() + "/" + c.DBFileName
}

func getDBDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	return homeDir + "/" + c.DBDirName

}

func checkOrCreatDBFiles() {
	if err := os.Mkdir(getDBDirPath(), 0777); err != nil && !os.IsExist(err) {
		log.Panic(err)
	}

	if _, err := os.Stat(GetDBFilePath()); os.IsNotExist(err) {
		_, err := os.Create(GetDBFilePath())
		if err != nil {
			log.Panic(err)
		}
	}
}
