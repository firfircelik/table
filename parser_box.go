package table

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var validSeparationLine = regexp.MustCompile(`^(((-\s)+-?)|(_+))\s*$`)
var whiteSpaceOnly = regexp.MustCompile(`^(\s*)$`)

// Key of box table
type Key struct {
	Column, Row string
}

// String implements Stringer
func (k *Key) String() string {
	return fmt.Sprintf("(%s,%s)", k.Column, k.Row)
}

// ParseBoxes source containing input inside boxes
func ParseBoxes(lines []string, columns int) (map[Key]string, error) {
	tableLines := limitToTable(lines)
	if len(tableLines) == 0 {
		return nil, errors.New("can't extract box")
	}
	entries, err := parseTable(tableLines, columns)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, errors.New("can't find any entries")
	}
	header := entries[0]

	result := map[Key]string{}
	for _, e := range entries[1:] {
		rowHeader := e[0]
		for i, rowEntry := range e[1:] {
			result[Key{header[i+1], rowHeader}] = rowEntry
		}
	}
	return result, nil
}

type entry []string

func (e entry) isNonEmpty() bool {
	for _, s := range e {
		if s != "" {
			return true
		}
	}
	return false
}

func parseTable(table []string, columns int) ([]entry, error) {
	var current entry
	var result []entry

	appendEntry := func(e entry) {
		if e.isNonEmpty() {
			result = append(result, e)
		}
	}

	for _, line := range table {
		if isSeparationLine(line) {
			appendEntry(current)
			current = make(entry, columns)
			continue
		}
		splitted := strings.Split(line, "|")
		if len(splitted) != columns &&
			(!onlyFirstColumnHasContent(splitted) || len(splitted) > columns) {

			return nil, errors.Errorf("unexpected number of columns %d, needed %d in %s",
				len(splitted), columns, line)
		}
		for i, elem := range splitted {
			current[i] = join(current[i], strings.Trim(elem, " "))
		}
	}
	appendEntry(current)
	return result, nil
}

// onlyFirstColumnHasContent returns true iff all but first columns are empty (first can be
// empty or not)
func onlyFirstColumnHasContent(row []string) bool {
	if len(row) <= 1 {
		return true
	}
	for _, elem := range row[1:] {
		if !isWhiteSpace(elem) {
			return false
		}
	}
	return true
}

func join(left, right string) string {
	right = strings.TrimSpace(right)
	if left == "" {
		return right
	}
	return left + " " + right
}

func limitToTable(lines []string) []string {
	var result []string
	boxStarted := false
	for _, line := range lines {
		boxStarted = boxStarted || isSeparationLine(line)
		if boxStarted {
			if !isInsideBox(line) {
				return result
			}
			if !isWhiteSpace(line) {
				result = append(result, line)
			}
		}
	}
	return result
}

func isWhiteSpace(line string) bool {
	return whiteSpaceOnly.MatchString(line)
}

func isInsideBox(line string) bool {
	return strings.ContainsAny(line, "|-_") || isWhiteSpace(line)
}

func isSeparationLine(line string) bool {
	return validSeparationLine.Match([]byte(line))
}
