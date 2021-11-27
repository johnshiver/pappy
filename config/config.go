package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"os"
)


const envPrefix = "PAPPY"
var c Config

type Config struct {
	DBDirName  string `required:"true" split_words:"true" default:".pappy"`
	DBFileName string `required:"true" split_words:"true" default:"data.db"`
}

func init() {
	envconfig.MustProcess(envPrefix, &c)
}

func GetDB() *sqlx.DB {
	checkOrCreatDBFiles()
	db, err := sqlx.Open("sqlite3", GetDBFilePath())
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func GetDBFilePath() string {
	return fmt.Sprintf("%s/%s", getDBDirPath(), c.DBFileName)
}

func getDBDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s/%s", homeDir, c.DBDirName)

}

func checkOrCreatDBFiles() {
	if err := os.Mkdir(getDBDirPath(), 0777); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	if _, err := os.Stat(GetDBFilePath()); os.IsNotExist(err) {
		_, err = os.Create(GetDBFilePath())
		if err != nil {
			log.Fatal(err)
		}
	}
}
