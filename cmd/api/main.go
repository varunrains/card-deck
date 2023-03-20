package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/varunrains/carddeck/internal/repository"
	"github.com/varunrains/carddeck/internal/repository/dbrepo"
)

const port = 8080

type application struct {
	DSN string
	DB  repository.DatabaseRepo
}

func main() {

	//set application config
	var app application

	//read from command line
	flag.StringVar(&app.DSN, "dsn", "host=postgres port=5432 user=postgres password=postgres dbname=cardsDB sslmode=disable timezone=UTC connect_timeout=5", "Postgres Connectionstring")
	flag.Parse()

	//connect to the database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
