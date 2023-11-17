package main

import (
	"api/config"
	"database/sql"
	"testing"

	"github.com/lib/pq"
)

func TestDB(t *testing.T) {

	pgUrl, err := pq.ParseURL(config.Postgrestring)
	db, err := sql.Open("postgres", pgUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

}
