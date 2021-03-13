package ormatic

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/saromanov/ormatic/generate"
)

// Ormatic defines the main structure
type Ormatic struct {
	db *sql.DB
}

// New provides initialization of the Ormatic
func New(db *sql.DB) (*Ormatic, error) {
	if db == nil {
		return nil, errors.New("db is not initialized")
	}
	return &Ormatic{
		db: db,
	}, nil
}

// Create provides creating of tables from structs
func (o *Ormatic) Create(d ...interface{}) error {
	return o.create(d...)
}

// Save provides saving of the data
func (o *Ormatic) Save(d interface{}) error {
	return o.save(d)
}

func (o *Ormatic) save(d interface{}) error {
	fields, err := getFieldsFromStruct(d)
	if err != nil {
		return errors.Wrap(err, "unable to get fields from the struct")
	}

	tableName := getObjectName(d)
	query, values, err := generate.Insert(tableName, fields)
	if err != nil {
		return errors.Wrap(err, "unable to generate statement")
	}
	_, err = o.db.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "unable to insert data")
	}
	return nil
}

func (o *Ormatic) create(models ...interface{}) error {
	if len(models) == 0 {
		return nil
	}
	for _, m := range models {
		getStructFieldsTypes(m)
		fields, err := getFieldsFromStruct(m)
		if err != nil {
			return errors.Wrap(err, "unable to get fields from the struct")
		}

		tableName := getObjectName(m)
		fmt.Println(tableName, fields)
	}
	return nil
}
