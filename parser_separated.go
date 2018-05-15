package table

import (
	"strings"

	"github.com/pkg/errors"
)

// ParseSeparated parses table assuming that each field is separated by two spaces
func ParseSeparated(lines []string, nbColumn int) (Parsed, error) {
	result := make([]parsedLine, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		splitted := Fields(line)
		if len(splitted) > nbColumn {
			return result, errors.Errorf(
				"can't parse table: too many columns: expected %d, got %d",
				len(splitted), nbColumn)
		}
		for len(splitted) <= nbColumn {
			splitted = append(splitted, "")
		}
		result[i] = parsedLine{parsed: splitted, original: line}
	}
	return result, nil
}
