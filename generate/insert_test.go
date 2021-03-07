package generate

import (
	"testing"

	"github.com/saromanov/ormatic/models"
	"github.com/stretchr/testify/assert"
)
func TestInsertBasic(t *testing.T) {
	_, err := Insert("", []models.Pair{})
	assert.Error(t, errNoTableName, err)
}
