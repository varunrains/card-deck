package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var counts int64

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) connectToDB() (*sql.DB, error) {
	for {
		connection, err := openDB(app.DSN)

		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println(("Connected to Postgres!"))
			return connection, nil
		}

		if counts > 10 {
			log.Println(err)
			return nil, err
		}

		log.Println("Backing off for two seconds...")

		time.Sleep(2 * time.Second)
		continue
	}
}
