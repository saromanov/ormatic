package ormatic

import (
	"database/sql"
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
	return &Ormatic {
		db: db,
	}, nil
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
	stat, err := generate.Insert("", fields)
	if err != nil {
		return errors.Wrap(err, "unable to generate statement")
	}
	fmt.Println(stat)
	return nil
}