package table

import (
	"math"

	"github.com/pkg/errors"
)

// column as discovered in the input text. It contains text start at from (inclusive) and ending
// at to (exclusive)
type column struct {
	from, to int
}

func (c *column) intersects(from, to int) bool {
	return c.to > from && to > c.from
}

func (c *column) contains(from, to int) bool {
	return c.from <= from && to <= c.to
}

func (c *column) extend(from, to int) {
	if from < c.from {
		c.from = from
	}
	if to > c.to {
		c.to = to
	}
}

// columns position is guessed from minimal left position and maximal right position
// among all rows
// nolint: splint, gocyclo
func columns(lines []string, nbColumn int) ([]column, error) {
	columns := make([]column, nbColumn)
	for i := range columns {
		// make sure that from and to will be updated on first row
		columns[i] = column{from: math.MaxInt32, to: 0}
	}
	foundAtLeastOneProperLine := false
	linesWithWrongLength := [][]field{}
	for _, line := range lines {
		fields := extractFields(line)
		if len(fields) != nbColumn {
			linesWithWrongLength = append(linesWithWrongLength, fields)
			continue
		}
		foundAtLeastOneProperLine = true
		for i, field := range fields {
			if field.from < columns[i].from {
				columns[i].from = field.from
			}
			if field.To() > columns[i].to {
				columns[i].to = field.To()
			}
		}
	}
	if !foundAtLeastOneProperLine {
		return columns, errors.Errorf("can't find any line with %d columns", nbColumn)
	}
	if columnsOverlap(columns) {
		return columns, errors.New(
			"columns in the table overlap, can't extract data in a reliable way")
	}
	for _, lineWithWrongLength := range linesWithWrongLength {
		for _, field := range lineWithWrongLength {
			for i, c := range columns {
				if !c.intersects(field.from, field.To()) || c.contains(field.from, field.To()) {
					continue
				}
				columnsCopy := append([]column{}, columns...)
				columnsCopy[i].extend(field.from, field.To())
				if !columnsOverlap(columnsCopy) {
					columns = columnsCopy
					break
				}
			}
		}
	}
	return columns, nil
}

func columnsOverlap(columns []column) bool {
	for i := 0; i < len(columns)-1; i++ {
		if columns[i].to > columns[i+1].from {
			return true
		}
	}
	return false
}
