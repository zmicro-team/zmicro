package form

import (
	"strings"
	"unicode"
)

// tagOptions is the slice of tag options in a struct field's "map"
// tag.
type tagOptions []string

// parseTag splits a struct field's tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	res := strings.Split(tag, ",")
	return res[0], res[1:]
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	for _, s := range o {
		if optionName == s {
			return true
		}
	}
	return false
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}
