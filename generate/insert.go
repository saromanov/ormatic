package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/saromanov/ormatic/models"
)

var (
	errNoTableName = errors.New("table name is not defined")
	errNoValues    = errors.New("values is not defined")
)

// Insert provides generation of Insert
func Insert(value models.Insert) (string, []interface{}, error) {
	if value.TableName == "" {
		return "", nil, errNoTableName
	}
	expr := ""
	expr += fmt.Sprintf("INSERT INTO %s (", value.TableName)
	num := len(value.Pairs)
	keys := make([]string, num)
	nums := make([]string, num)
	data := make([]interface{}, num)
	for i, v := range value.Pairs {
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
