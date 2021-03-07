package ormatic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldsFromStruct(t *testing.T) {
	_, err := getFieldsFromStruct(nil)
	assert.Error(t, errNoStruct, err)
	_, err = getFieldsFromStruct(4)
	assert.Error(t, errNoStruct, err)
	_, err = getFieldsFromStruct([]string{})
	assert.Error(t, errNoStruct, err)
}
