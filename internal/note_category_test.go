package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoteCategoryString(t *testing.T) {
	assert.Equal(t, "top", TopNote.String())
	assert.Equal(t, "middle", MiddleNote.String())
	assert.Equal(t, "base", BaseNote.String())
	assert.Equal(t, "uncategorized", UncategorizedNote.String())
	assert.Equal(t, "", NoteCategory(-1).String())
}

func TestNoteCategoryFromString(t *testing.T) {
	top, err := NoteCategoryFromString("Top")
	assert.Nil(t, err)
	assert.Equal(t, TopNote, top)

	middle, err := NoteCategoryFromString("Middle")
	assert.Nil(t, err)
	assert.Equal(t, MiddleNote, middle)

	base, err := NoteCategoryFromString("Base")
	assert.Nil(t, err)
	assert.Equal(t, BaseNote, base)

	uncategorized, err := NoteCategoryFromString("Uncategorized")
	assert.Nil(t, err)
	assert.Equal(t, UncategorizedNote, uncategorized)

	unknown, err := NoteCategoryFromString("foo")
	assert.Error(t, err, "foo")
	assert.Equal(t, NoteCategory(-1), unknown)
}
