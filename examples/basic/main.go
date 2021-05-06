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
	Address Address `db:"address_id" orm:"ON=address.id"`
}

type Address struct {
	ID int `orm:"PRIMARY_KEY,NOT_NULL"`
	Name  string
	Title string
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

	/*if err := o.Create(&Book{}); err != nil {
		panic(err)
	}*/

	if err := o.Save(&Book{
		ID:    8,
		Title: "test",
		Address: Address{
			ID: 23,
			Name: "Moskvaa",
			Title: "builinga",
		},
	}); err != nil {
		panic(err)
	}

	resp, err := o.Find(&Book{
		Title: "test",
	}).Join(ormatic.DefaultJoin, Address{}).Many(Book{})
	if err != nil {
		panic(err)
	}
	fmt.Println("OUTPUT: ", resp)

	respOne, err := o.Find(&Book{
		Title: "test",
	}).One(Book{})
	if err != nil {
		panic(err)
	}
	fmt.Println("OUTPUT: ", respOne.(Book))
}
