package table

import (
	"strings"
)

type field struct {
	value string
	from  int
}

func (f *field) To() int { return f.from + len(f.value) }

func extractFields(s string) []field {
	var linePosition int
	fieldValues := Fields(s)
	fields := make([]field, 0, len(fieldValues))
	for _, fieldValue := range fieldValues {
		// can happen for the first or last column
		if fieldValue == "" {
			continue
		}
		index := strings.Index(s, fieldValue)
		fields = append(fields, field{value: fieldValue, from: index + linePosition})
		s = s[index+len(fieldValue):]
		linePosition += index + len(fieldValue)
	}
	return fields
}

// Fields returns string splitted into fields. Each field is separated by two or more
// whitespaces
func Fields(s string) []string { return twoOrMoreWhitespaces.Split(s, -1) }
