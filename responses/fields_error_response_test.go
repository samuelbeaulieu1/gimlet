package responses

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFieldsError(t *testing.T) {
	err := NewFieldsError([]string{}, []string{})

	assert.NotNil(t, err)
	assert.Equal(t, err.Status, "fail")
}

func TestFieldsErrorMsg(t *testing.T) {
	errMsg := []string{"testfield1 err", "testfield2 err"}
	errFields := []string{"test1", "test2"}
	err := NewFieldsError(errMsg, errFields)

	assert.NotNil(t, err)
	assert.Equal(t, err.Data.Fields, errFields)
	assert.Equal(t, err.Data.Messages, errMsg)

	message := fmt.Sprintf("%s\n%s", errMsg[0], errMsg[1])
	assert.Equal(t, err.Error(), message)
}
