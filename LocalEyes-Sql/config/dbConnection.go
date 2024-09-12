package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var dbClient *sql.DB
var once sync.Once

func GetSQLClient() *sql.DB {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("DBUser"),
			os.Getenv("DBPassword"),
			os.Getenv("DBHost"),
			os.Getenv("DBPort"),
			os.Getenv("DBName"),
		)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}

		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}

		dbClient = db
	})
	dbClient.SetConnMaxLifetime(time.Hour * 1)
	return dbClient
}

func CloseDBClient() {
	err := dbClient.Close()
	if err != nil {
		fmt.Println("Error closing DB client")
		return
	}
}
