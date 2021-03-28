package ormatic

import (
	"database/sql"
	"fmt"
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
		data = append(data, fmt.Sprintf("%s=%s ", f.Key, f.Value))
	}
	stat += strings.Join(data, "AND")
	return stat, nil
}
