package ormatic

import "github.com/saromanov/ormatic/models"

// FindResult returns
type FindResult struct {
	limit          uint
	table          string
	err            error
	nonEmptyFields []models.Pair
	selectedFields []models.Pair
	fields         []models.Pair
}

// Do returns result of the query
func (d *FindResult) Do(dest interface{}) error {
	if d.err != nil {
		return d.err
	}
	return nil
}
