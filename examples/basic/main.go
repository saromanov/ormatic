package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/saromanov/ormatic"
)

type Book struct {
	Title string
	Address Address
}

type Address struct {
	Name string
	Basic Another
}

type Another struct {
	Basic string
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		"tracer", "tracer", "tracer")
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	o, err := ormatic.New(db)
	if err != nil {
		panic(err)
	}

	if err := o.Create(&Book{}); err != nil {
		panic(err)
	}
}
