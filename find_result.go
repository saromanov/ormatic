package ormatic

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
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
	properies      FindProperties
	joins          []string
}

// FindProperties defines properties for find
type FindProperties struct {
	limit   uint
	orderBy []string
	or      []models.Pair
}

// One returns single result of the query
func (d *FindResult) One(m interface{}) (interface{}, error) {
	d.limit = 1
	result, err := d.Many(m)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find data")
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result[0], nil
}

// Many returns multiple result of the query
func (d *FindResult) Many(m interface{}) ([]interface{}, error) {
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
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("unable to close rows: ", err)
		}
	}()

	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get number of columns")
	}
	fmt.Println("COLUMS: ", columns)
	resp := []interface{}{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, errors.Wrap(err, "unable to scan value")
		}

		row := map[string]interface{}{}
		for i, v := range values {
			row[columns[i]] = v
		}
		if err := mapstructure.Decode(row, &m); err != nil {
			return nil, errors.Wrap(err, "unable to decode result to struct")
		}
		resp = append(resp, m)
	}
	_, err = d.db.Exec(res)
	if err != nil {
		return nil, errors.Wrap(err, "unable to execute statement")
	}
	return resp, nil
}

// Limit sets limit of results from the query
func (d *FindResult) Limit(limit uint) *FindResult {
	d.properies.limit = limit
	return d
}

// OrderBy sets params for sorring
func (d *FindResult) OrderBy(params []string) *FindResult {
	d.properies.orderBy = params
	return d
}

// Join provides joining of the tables
func (d *FindResult) Join(joinType string, t interface{}) *FindResult {
	d.joins = append(d.joins, getObjectName(t))
	return d
}

// Or sets or at the where statement
func (d *FindResult) Or(params map[string]interface{}) *FindResult {
	result := []models.Pair{}
	for key, value := range params {
		result = append(result, models.Pair{
			Key:   key,
			Value: value,
		})
	}
	d.properies.or = result
	return d
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
		if f.Join.Source != "" {
			continue
		}
		data = append(data, fmt.Sprintf("%s=%s ", f.Key, d.setValue(f.Value)))
	}
	stat += strings.Join(data, "AND")
	if len(d.properies.or) > 0 {
		data := make([]string, 0, len(d.nonEmptyFields))
		for _, f := range d.nonEmptyFields {
			data = append(data, fmt.Sprintf("%s=%s ", f.Key, d.setValue(f.Value)))
		}
		stat += strings.Join(data, "OR")
	}
	if len(d.properies.orderBy) > 0 {
		stat += " ORDER BY " + strings.Join(d.properies.orderBy, ",")
	}
	if d.properies.limit != 0 {
		stat += fmt.Sprintf(" LIMIT %d", d.properies.limit)
	}
	fmt.Println("EXPR: ", stat)
	return stat, nil
}

// check if value is a string then add it with commas
func (d *FindResult) setValue(value interface{}) string {
	if reflect.ValueOf(value).Kind() == reflect.String {
		return fmt.Sprintf("'%s'", value)
	}
	return fmt.Sprintf("%v", value)
}
