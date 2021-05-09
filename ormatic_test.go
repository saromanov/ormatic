package ormatic

import (
	"database/sql"
	"fmt"
	"testing"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New(nil)
	assert.Error(t, err)

	db := newDB(t)
	_, err = New(db)
	assert.NoError(t, err)
}

func newDB(t *testing.T) *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		"tracer", "tracer", "tracer")
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
