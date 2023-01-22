package gimlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewID(t *testing.T) {
	id, err := CreateNewID()

	assert.Nil(t, err)
	assert.Greater(t, len(id), 0)
}
