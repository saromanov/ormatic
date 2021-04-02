package ormatic

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/saromanov/ormatic/models"
)

// FindResult returns
type FindResult struct {
	db             *sql.DB
	limit          uint
	table          string
	err            error
	nonEmptyFields []models.Pair
	selectedFields []models.Pair
	fields         []models.Pair
}

// Do returns result of the query
func (d *FindResult) Do(dest interface{}) error {
	if d.db == nil {
		return errors.New("db is not defined")
	}
	if d.err != nil {
		return d.err
	}
	res, err := d.constructFindStatement(d.table, d.nonEmptyFields)
	if err != nil {
		return err
	}
	fmt.Println(res)
	rows, err := d.db.Query(res)
	if err != nil {
		return errors.Wrap(err, "unable to query data")
	}
	defer func(){
		if err := rows.Close(); err != nil {
			log.Println("unable to close rows: ", err)
		}
	}()

	for rows.Next(){
		data := make([]interface{}, 2)
		for i, _ := range data {
			var res interface{}
			data[i] = &res
		}
		if err := rows.Scan(data...); err != nil {
			return errors.Wrap(err, "unable to scan value")
		}
		first := data[0]
		fmt.Println(first.(*interface{}))
	}
	_, err = d.db.Exec(res)
	if err != nil {
	  return errors.Wrap(err, "unable to execute statement")
	}
	return nil
}

// constructFindStatement provides constructing find statement
// like SELECT * FROM value1=foo AND value2=bar;
func (d *FindResult) constructFindStatement(tableName string, nonEmptyFields []models.Pair) (string, error) {
	stat := "SELECT * FROM " + tableName
	if len(nonEmptyFields) == 0 {
		return stat + ";", nil
	}
	stat += " WHERE "
	data := make([]string, 0, len(nonEmptyFields))
	for _, f := range nonEmptyFields {
		data = append(data, fmt.Sprintf("%s=%s ", f.Key, d.setValue(f.Value)))
	}
	stat += strings.Join(data, "AND")
	return stat, nil
}

// check if value is a string then add it with commas
func (d *FindResult) setValue(value interface{}) string {
	if reflect.ValueOf(value).Kind() == reflect.String {
		return fmt.Sprintf("'%s'", value)
	}
	return fmt.Sprintf("%v", value)
}