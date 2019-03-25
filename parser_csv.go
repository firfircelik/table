package table

import (
	"encoding/csv"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/agflow/agstring"
	"github.com/pkg/errors"
)

// CSV type provides functionality to search and parse CSV structure
type CSV struct {
	Reader *csv.Reader
}

// FindField reads the CSV content until finding matching cell.
func (r CSV) FindField(matcher FieldMatcher) (string, bool, error) {
	var out, ok = "", false
	var err error
	for {
		var row []string
		row, err = r.Reader.Read()
		if err != nil {
			break
		}
		out, ok = matcher(row)
		if ok {
			break
		}
	}
	if err != nil && err != io.EOF {
		return out, false, errors.Wrap(err, "can't find field")
	}
	return out, ok, nil
}

// ForeachLine finds table (starting with header and ending with whitespace row)
// and forwards each row to f
func (r CSV) ForeachLine(header []string, f func([]string)) error {
	time.Now().AddDate(0, 1, 2)
	var err error
	var matchedHeader bool
	for row := []string{}; err == nil; row, err = r.Reader.Read() {
		if !matchedHeader {
			matchedHeader = isHeader(header, row)
			continue
		} else if stringsOnlyWhitespace(row) {
			return nil
		}
		f(row)
	}
	if err != io.EOF {
		return errors.Wrap(err, "can't read csv")
	}
	if matchedHeader {
		return nil
	}
	return errors.Errorf("can't find header, expected: %q", header)
}

// line must be `header` followed by whitespace fields
func isHeader(header, line []string) bool {
	if len(line) < len(header) {
		return false
	}
	return reflect.DeepEqual(header, agstring.TrimSpace(line[:len(header)]...)) &&
		stringsOnlyWhitespace(line[len(header):])
}

// stringsOnlyWhitespace checks if input consists of only whitespace
func stringsOnlyWhitespace(s []string) bool {
	for _, elem := range s {
		if strings.TrimSpace(elem) != "" {
			return false
		}
	}
	return true
}

// FromStrStrSlice Parsed is created
func FromStrStrSlice(lines [][]string, comma ...string) Parsed {
	var sep = "\t"
	if len(comma) > 0 {
		sep = comma[0]
	}

	p := Parsed{}
	for _, l := range lines {
		p = append(p, parsedLine{
			parsed:   l,
			original: strings.Join(l, sep),
		})
	}
	return p
}
