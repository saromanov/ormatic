package generate

import (
	"fmt"
	"strings"
	"errors"

	"github.com/saromanov/ormatic/models"
)

var (
	errNoTableName = errors.New("table name is not defined")
)

// Insert provides generation of Insert
func Insert(tableName string, values []models.Pair) (string, error) {
	if tableName == "" {
		return "", errNoTableName
	}
	expr := fmt.Sprintf("INSERT INTO %s (", tableName)
	num := len(values)
	keys := make([]string, num)
	nums := make([]string, num)
	for i, v := range values {
		keys[i] = v.Key
		nums[i] = fmt.Sprintf("$%d", i+1)
	}
	keysStr := strings.Join(keys, ",")
	expr += keysStr + ") "
	expr += "VALUES (" + strings.Join(nums, ",") + ")"
	return expr, nil
}