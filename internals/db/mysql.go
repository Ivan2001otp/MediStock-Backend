package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var dbInstance *sql.DB

var once sync.Once


func Close() error {
	if dbInstance != nil {
		log.Println("Closing database connection pool.")
		return dbInstance.Close()
	}
	return nil;
}

func Connect(ctx context.Context) (*sql.DB, error) {
	var err error

	once.Do(
		func() {

			err = godotenv.Load()
			if err != nil {
				log.Println("Something went wrong. Could not load environment variables.")
				return
			}

			dbUser := os.Getenv("DB_USER")
			dbPassword := os.Getenv(("DB_PASSWORD"))
			dbHost := os.Getenv("DB_HOST")
			dbPort := os.Getenv("DB_PORT")
			dbName := os.Getenv("DB_NAME")

			if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
				err = fmt.Errorf("database environment variables (DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME) must be set")
				return
			}

			var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

			log.Println("Data source name is : ", dsn)

			dbInstance, err = sql.Open("mysql", dsn)

			if err != nil {
				err = fmt.Errorf("failed to open db connection : %w ", err)
				return
			}

			// set max connections
			dbInstance.SetMaxOpenConns(25)

			//set max idle connections
			dbInstance.SetMaxIdleConns(10)

			//set conn maxlife time
			dbInstance.SetConnMaxLifetime(5 * time.Minute)

			// ping the db connection to verify whether its alive or not
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

			defer cancel()

			err = dbInstance.PingContext(pingCtx)

			if err != nil {
				dbInstance.Close()
				dbInstance = nil
				err = fmt.Errorf("failed to connect to database : %w", err)
				return
			}

			log.Println("Successfully connected to MySql database and configured connection pool !")
		})

	if err != nil {
		return nil, err
	}

	return dbInstance, nil
}
