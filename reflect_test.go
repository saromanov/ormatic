package ormatic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Test struct {
	
}

func TestGetFieldsFromStruct(t *testing.T) {
	_, err := getFieldsFromStruct(nil)
	assert.Error(t, errNoStruct, err)
	_, err = getFieldsFromStruct(4)
	assert.Error(t, errNoStruct, err)
	_, err = getFieldsFromStruct([]string{})
	assert.Error(t, errNoStruct, err)
}

func TestGetStructName(t *testing.T) {
	assert.Equal(t, "test", getObjectName(&Test{}))
}
