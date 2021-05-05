package ormatic

import (
	"database/sql"
	"fmt"

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

// Delete provides deleteing of data
func (o *Ormatic) Delete(d interface{}) error {
	return o.delete(d)
}

func (o *Ormatic) Find(query interface{}) *FindResult {
	return o.find(query)
}

func (o *Ormatic) save(d interface{}) error {
	fields, err := prepareInsert(d)
	if err != nil {
		return errors.Wrap(err, "unable to get fields from the struct")
	}

	for _, v := range fields {
		query, values, err := generate.Insert(v)
		if err != nil {
			return errors.Wrap(err, "unable to generate statement")
		}
		_, err = o.db.Exec(query, values...)
		if err != nil {
			return errors.Wrap(err, "unable to insert data")
		}
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

func (o *Ormatic) find(q interface{}) *FindResult {
	fields, err := getFieldsFromStruct(q)
	res := &FindResult{
		table: getObjectName(q),
		db:    o.db,
	}
	if err != nil {
		res.err = errors.Wrap(err, "unable to get fields from the struct")
		return res
	}
	res.nonEmptyFields = fields
	return res
}

func (o *Ormatic) delete(d interface{}) error {
	_, err := getFieldsFromStruct(d)
	if err != nil {
		return errors.Wrap(err, "unable to get fields from the struct")
	}
	return nil
}

// consructCreateTable provides generation of the create
// table statement
func (o *Ormatic) constructCreateTable(models []models.Create) error {
	text := "BEGIN TRANSACTION;\n"
	for _, m := range models {
		text += fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", m.TableName)
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
		text += ");"
		if len(m.Relationships) > 0 {
			/*for _, r := range m.Relationships {
				constraint := fmt.Sprintf("fk_%s%s%s", m.TableName, r.TableName, "test")
				text += fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(%s);", m.TableName, constraint, r.Name, r.TableName, r.Column)
			}*/
		}
		text += "\n"
	}
	text += "\nCOMMIT;"
	fmt.Println("TEXT: ", text)
	if _, err := o.exec(text); err != nil {
		return errors.Wrap(err, "unable to execute data")
	}
	return nil
}

func (o *Ormatic) exec(query string) (sql.Result, error) {
	return o.db.Exec(query)
}
