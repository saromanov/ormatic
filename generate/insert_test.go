package generate

import (
	"testing"

	"github.com/saromanov/ormatic/models"
	"github.com/stretchr/testify/assert"
)
func TestInsertBasic(t *testing.T) {
	_, _, err := Insert("", []models.Pair{})
	assert.Error(t, errNoTableName, err)

	data, values, err := Insert("test", []models.Pair{models.Pair{
		Key: "key",
		Value: "value",
	}})
	assert.NoError(t, err)
	assert.Equal(t, data, "INSERT INTO test (key) VALUES ($1)")
	assert.Equal(t, "value", values[0].(string))

	data, values, err = Insert("test2", []models.Pair{models.Pair{
		Key: "key",
		Value: "value",
	},

	models.Pair{
		Key:"foo",
		Value:"bar",
	}})
	assert.NoError(t, err)
	assert.Equal(t, data, "INSERT INTO test2 (key,foo) VALUES ($1,$2)")
	assert.Equal(t, "value", values[0].(string))
	assert.Equal(t, "bar", values[1].(string))
}
