package testutil

import (
	"errors"
	"regexp"
)

const uuidRegexp = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`

var uuidCompiledRegexp = regexp.MustCompile(uuidRegexp)

var errNoUUIDFound = errors.New("no uuid found")
var errEmptyText = errors.New("empty text")

func ExtractFirstUUIDFromText(text string) (string, error) {
	if text == "" {
		return "", errEmptyText
	}

	match := uuidCompiledRegexp.FindString(text)
	if match == "" {
		return "", errNoUUIDFound
	}

	return match, nil
}
