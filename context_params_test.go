package gimlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContextParams(t *testing.T) {
	keys := []string{"test"}
	values := []string{"value"}

	contextParams := NewContextParams(keys, values)

	assert.NotNil(t, contextParams)
	assert.Equal(t, len(keys), len(contextParams.params))
}

func TestNewContextParamsGet(t *testing.T) {
	contextParams := ContextParams{
		params: map[string]string{"test": "value"},
	}

	value := contextParams.Get("test")
	assert.Equal(t, value, "value")
}

func TestNewContextParamsGetEmpty(t *testing.T) {
	contextParams := ContextParams{
		params: map[string]string{},
	}

	value := contextParams.Get("test")
	assert.Equal(t, value, "")
}
