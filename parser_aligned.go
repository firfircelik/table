package table

import (
	"unicode/utf8"

	"github.com/pkg/errors"
)

func getAlignedOptions(options ...*Options) (*Options, error) {
	opts := getOptions(options...)
	if len(opts.NumberOfColumns) == 0 {
		return opts, errors.New("at least one column count should be given")
	}
	return opts, nil
}

// ParseAligned table. Tries to parse table which looks like this:
// a   b   c
// aa  bb  c
// a  bbb  c
//
// Please consult tests to see in which situations this function
// performs well and in which it does not.
// Number of columns in the table needs to be provided.
// If none of the rows has the expected number of column, error is returned.
// WARNING: This function does work only with ASCII strings
// (so result might be wrong for UTF-8 strings)
func ParseAligned(lines []string, options ...*Options) (Parsed, error) {
	opts, err := getAlignedOptions(options...)
	if err != nil {
		return Parsed{}, err
	}
	for _, c := range opts.NumberOfColumns {
		if p, err := parseAligned(lines, c); err == nil {
			return p, nil
		}
	}
	return Parsed{}, errors.Errorf("can't parse lines with column counts: %v",
		opts.NumberOfColumns)
}

func parseAligned(lines []string, nbColumn int) (Parsed, error) {
	cols, err := columns(lines, nbColumn)
	if err != nil {
		return nil, err
	}
	result := make([]parsedLine, len(lines))
	for i, line := range lines {
		result[i] = parsedLine{
			parsed:   splitByCols(line, cols),
			original: line}
	}
	return result, nil
}

func splitByCols(line string, cols []column) []string {
	splitted := make([]string, len(cols))
	for i, c := range cols {
		if len(line) >= c.to {
			from := findFrom(line, cols, i)
			to := findTo(line, cols, i)
			splitted[i] = line[from:to]
		} else if len(line) > c.from {
			from := findFrom(line, cols, i)
			splitted[i] = line[from:]
		} // if none of those two match splitted[i] will contains empty string
	}
	return splitted
}

func findFrom(line string, cols []column, i int) int {
	prevColumn := 0
	if i > 0 {
		prevColumn = cols[i-1].to
	}
	from := cols[i].from
	for {
		_, w := utf8.DecodeLastRuneInString(line[:from])
		// we can't extend return current
		if w == 0 {
			return from
		}
		// we have reached the line start
		if from-w <= 0 {
			return 0
		}
		// by extending the column we overlapped the previous column, lets return the original from
		if from-w <= prevColumn {
			return from
		}
		from -= w
	}
}

func findTo(line string, cols []column, i int) int {
	nextColumn := len(line)
	if i < len(cols)-1 {
		nextColumn = cols[i+1].from
	}
	to := cols[i].to
	for {
		_, w := utf8.DecodeRuneInString(line[to:])
		// we can't extend return current
		if w == 0 {
			return to
		}
		// we have reached the line end, we are in the last column
		if to+w == len(line) {
			return len(line)
		}
		// by extending the column we overlapped the next column, lets return the original from
		if to+w >= nextColumn {
			return to
		}
		to += w
	}
}

// DownFrom generates decreasing column numbers from start to end
func DownFrom(start uint, end ...uint) []int {
	s, e := int(start), 0
	if len(end) > 0 {
		e = int(end[0])
	}
	cols := make([]int, s-e)
	for i := range cols {
		cols[i] = s - i
	}
	return cols
}
