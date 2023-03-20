package main

import (
	"os"
	"testing"

	"github.com/varunrains/carddeck/internal/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {

	app.DB = &dbrepo.TestDBRepo{}
	os.Exit(m.Run())
}
