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
		var data interface{}
		var data2 interface{}
		if err := rows.Scan(&data, &data2); err != nil {
			return errors.Wrap(err, "unable to scan value")
		}
		fmt.Println(data)
	}
	_, err = d.db.Exec(res)
	if err != nil {
	  return errors.Wrap(err, "unable to execute statement")
	}
	return nil
}

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

func (d *FindResult) setValue(value interface{}) string {
	if reflect.ValueOf(value).Kind() == reflect.String {
		return fmt.Sprintf("'%v'", value)
	}
	return fmt.Sprintf("%v", value)
}