package internal

import (
	"fmt"
	"strings"
)

type NoteCategory int

const (
	TopNote NoteCategory = iota
	MiddleNote
	BaseNote
	UncategorizedNote
)

var NoteCategoryMap = map[string]NoteCategory{
	"top":           TopNote,
	"middle":        MiddleNote,
	"base":          BaseNote,
	"uncategorized": UncategorizedNote,
}

func NoteCategoryFromString(s string) (NoteCategory, error) {
	noteCategory, ok := NoteCategoryMap[strings.ToLower(s)]
	if !ok {
		return -1, fmt.Errorf("unknown note category: %s", s)
	}

	return noteCategory, nil
}

func (c NoteCategory) String() string {
	switch c {
	case TopNote:
		return "top"
	case MiddleNote:
		return "middle"
	case BaseNote:
		return "base"
	case UncategorizedNote:
		return "uncategorized"
	default:
		return ""
	}
}
