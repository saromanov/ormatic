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
	properies FindProperties
}

// FindProperties defines properties for find
type FindProperties struct {
	limit uint
	orderBy string 
}

// Do returns result of the query
func (d *FindResult) Do() ([][]interface{}, error) {
	if d.db == nil {
		return nil, errors.New("db is not defined")
	}
	if d.err != nil {
		return nil, d.err
	}
	res, err := d.constructFindStatement()
	if err != nil {
		return nil, err
	}
	rows, err := d.db.Query(res)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query data")
	}
	defer func(){
		if err := rows.Close(); err != nil {
			log.Println("unable to close rows: ", err)
		}
	}()
	
	columns, err := rows.Columns()
	if err != nil {
	    return nil, errors.Wrap(err, "unable to get number of columns") 
	}
	resp := [][]interface{}{}
	for rows.Next(){
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
            valuePtrs[i] = &values[i]
        }
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, errors.Wrap(err, "unable to scan value")
		}
		resp = append(resp, values)
	}
	_, err = d.db.Exec(res)
	if err != nil {
	  return nil, errors.Wrap(err, "unable to execute statement")
	}
	return resp, nil
}

// constructFindStatement provides constructing find statement
// like SELECT * FROM value1=foo AND value2=bar;
func (d *FindResult) constructFindStatement() (string, error) {
	stat := "SELECT * FROM " + d.table
	if len(d.nonEmptyFields) == 0 {
		return stat + ";", nil
	}
	stat += " WHERE "
	data := make([]string, 0, len(d.nonEmptyFields))
	for _, f := range d.nonEmptyFields {
		data = append(data, fmt.Sprintf("%s=%s ", f.Key, d.setValue(f.Value)))
	}
	stat += strings.Join(data, "AND")
	if d.properies.orderBy != "" {
		stat += " ORDER BY " + d.properies.orderBy
	}
	if d.properies.limit != 0 {
		stat += fmt.Sprintf(" LIMIT %d", d.properies.limit)
	}
	return stat, nil
}

// check if value is a string then add it with commas
func (d *FindResult) setValue(value interface{}) string {
	if reflect.ValueOf(value).Kind() == reflect.String {
		return fmt.Sprintf("'%s'", value)
	}
	return fmt.Sprintf("%v", value)
}