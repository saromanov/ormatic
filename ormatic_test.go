package ormatic

import (
	"database/sql"
	"fmt"
	"testing"
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

type Test1 struct {
	ID      int `orm:"PRIMARY_KEY,NOT_NULL"`
	Title   string
}

func TestNew(t *testing.T) {
	_, err := New(nil)
	assert.Error(t, err)

	db := newDB(t)
	_, err = New(db)
	defer db.Close()
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	db := newDB(t)
	defer db.Close()
	orm, err := New(newDB(t))
	assert.NoError(t, err)
	assert.NoError(t, orm.Create(&Test1{}))
}

func TestInsert(t *testing.T) {
	db := newDB(t)
	dropTable(t, db, "test1")
	defer db.Close()
	orm, err := New(newDB(t))
	assert.NoError(t, err)
	assert.NoError(t, orm.Create(&Test1{}))
	assert.NoError(t, orm.Save(&Test1{
		ID: 1,
		Title:"test",
	}))
	dropTable(t, db, "test1")
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

func dropTable(t *testing.T, db *sql.DB, name string) {
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", name))
	if err != nil {
		t.Fatal(err)
	}
}