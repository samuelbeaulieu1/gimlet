package responses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	err := NewError("")

	assert.NotNil(t, err)
	assert.Equal(t, err.Status, "error")
}

func TestErrorMsg(t *testing.T) {
	errMsg := "error msg"
	err := NewError(errMsg)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), errMsg)
}
