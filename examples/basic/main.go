package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/saromanov/ormatic"
)

type Book struct {
	ID      int `orm:"PRIMARY_KEY,NOT_NULL"`
	Title   string
	Address Address
}

type Address struct {
	Name   string
	Basic  Another
	BookID int `orm:"ON=book.id"`
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

	/*if err := o.Save(&Book{
		ID:    15,
		Title: "test",
	}); err != nil {
		panic(err)
	}*/

	var books []Book
	err = o.Find(&Book{
		Title: "test",
	}).Do(&books)
	if err != nil {
		panic(err)
	}
}
