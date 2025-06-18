package testutil

import (
	"errors"
	"regexp"
)

const uuidRegExp = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`

var errNoUUIDFound = errors.New("no uuid found")
var errEmptyText = errors.New("empty text")

func ExtractFirstUUIDFromText(text string) (string, error) {
	if text == "" {
		return "", errEmptyText
	}

	re := regexp.MustCompile(uuidRegExp)
	match := re.FindString(text)
	if match == "" {
		return "", errNoUUIDFound
	}

	return match, nil
}
