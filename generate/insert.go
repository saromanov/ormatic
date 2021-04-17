package generate

import (
	"fmt"
	"strings"
	"errors"

	"github.com/saromanov/ormatic/models"
)

var (
	errNoTableName = errors.New("table name is not defined")
	errNoValues = errors.New("values is not defined")
)

// Insert provides generation of Insert
func Insert(tableName string, values []models.Pair) (string, []interface{}, error) {
	if tableName == "" {
		return "", nil, errNoTableName
	}
	expr := fmt.Sprintf("INSERT INTO %s (", tableName)
	num := len(values)
	keys := make([]string, num)
	nums := make([]string, num)
	data := make([]interface{}, num)
	for i, v := range values {
		keys[i] = v.Key
		data[i] = v.Value
		nums[i] = fmt.Sprintf("$%d", i+1)
	}
	keysStr := strings.Join(keys, ",")
	expr += keysStr + ") "
	if len(nums) == 0 {
		return "", nil, errNoValues
	}
	expr += "VALUES (" + strings.Join(nums, ",") + ")"
	return expr, data, nil
}