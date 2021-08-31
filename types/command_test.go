package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandNotEmpty(t *testing.T) {
	com := Command()

	assert.IsType(t, command{}, com)

	v := reflect.ValueOf(com)
	for i := 0; i < v.NumField(); i++ {
		assert.NotEmpty(t, v.Field(i).Interface())
	}
}
