package ormatic

import (
	"database/sql"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/saromanov/ormatic/generate"
	"github.com/saromanov/ormatic/models"
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

// Drop provides drop of the table
func (o *Ormatic) Drop(table string) error {
	return o.drop(table)
}

// Driver returns sql.DB driver
func (o *Ormatic) Driver() *sql.DB {
	return o.db
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

func (o *Ormatic) drop(table string) error {
	if table == "" {
		return nil
	}
	_, err := o.exec(fmt.Sprintf("DROP TABLE %s", table))
	if err != nil {
		return errors.Wrap(err, "unable to drop tablle")
	}
	return nil
}

func (o *Ormatic) create(models ...interface{}) error {
	if len(models) == 0 {
		return nil
	}
	for _, m := range models {
		fields, err := getStructFieldsTypes(m)
		if err != nil {
			return errors.Wrap(err, "unable to get fields from the struct")
		}
		if err := o.constructCreateTable(fields); err != nil {
			return errors.Wrap(err, "unable to execute create table")
		}
	}
	return nil
}

// consructCreateTable provides generation of the create
// table statement
func (o *Ormatic) constructCreateTable(models []models.Create) error {
	for _, m := range models {
		text := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", m.TableName)
		if len(m.TableFields) == 0 {
			_, err := o.exec(text)
			if err != nil {
				return errors.Wrap(err, "unable to execute data")
			}
			continue
		}
		text += "("
		for i, f := range m.TableFields {
			text += fmt.Sprintf("%s %s", f.Name, f.Type)
			if f.Tags.PrimaryKey {
				text += " PRIMARY KEY"
			}
			if f.Tags.NotNULL {
				text += " NOT NULL"
			}
			if f.Tags.Unique {
				text += " UNIQUE"
			}
			if (len(m.TableFields) - i) != 1 {
				text += ","
			}
		}
		text += ")"
		if _, err := o.exec(text); err != nil {
			return errors.Wrap(err, "unable to execute data")
		}
	}
	return nil
}

func (o *Ormatic) exec(query string) (sql.Result, error) {
	return o.db.Exec(query)
}
