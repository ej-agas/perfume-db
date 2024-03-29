package nanoid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"
	length := 12

	generator := NewNanoIdGenerator(alphabet, length)
	id, err := generator.Generate()

	assert.Nil(t, err)
	assert.Equal(t, length, len(id))
}
