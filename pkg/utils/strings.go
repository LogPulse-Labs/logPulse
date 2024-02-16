package utils

import (
	"regexp"
	"strings"
)

// CopyStrings makes a deep copy of the passed in string slice and returns
// the copy.
func CopyStrings(in []string) []string {
	if in == nil {
		return nil
	}

	out := make([]string, len(in))
	copy(out, in)

	return out
}

func ToSnakeCase(s string) string {
	s = strings.ReplaceAll(s, " ", "")

	var (
		re     = regexp.MustCompile("(.)([A-Z][a-z]+)")
		re2    = regexp.MustCompile("([a-z0-9])([A-Z])")
		result = re.ReplaceAllString(s, "${1}_${2}")
	)

	result = re2.ReplaceAllString(result, "${1}_${2}")
	return strings.ToLower(result)
}
