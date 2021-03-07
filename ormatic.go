package ormatic

import (
	"database/sql"
	"github.com/pkg/errors"
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
	return nil
}