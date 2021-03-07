package generate

import (
	"testing"

	"github.com/saromanov/ormatic/models"
	"github.com/stretchr/testify/assert"
)
func TestInsertBasic(t *testing.T) {
	_, err := Insert("", []models.Pair{})
	assert.Error(t, errNoTableName, err)

	data, err := Insert("test", []models.Pair{models.Pair{
		Key: "key",
		Value: "value",
	}})
	assert.NoError(t, err)
	assert.Equal(t, data, "INSERT INTO test (key) VALUES ($1)")
}
