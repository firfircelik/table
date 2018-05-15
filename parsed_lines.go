package table

// Parsed represents parsed aligned table
type Parsed []parsedLine

type parsedLine struct {
	original string
	parsed   []string
}

// FindLine returns parsed version of the first line matching predicate
func (p Parsed) FindLine(predicate func(string) bool) []string {
	for _, line := range p {
		if predicate(line.original) {
			return line.parsed
		}
	}
	return nil
}

// Lines returns parsed line
func (p Parsed) Lines() [][]string {
	result := make([][]string, len(p))
	for i, line := range p {
		result[i] = line.parsed
	}
	return result
}

// Head returns first parsed line
func (p Parsed) Head() ([]string, bool) {
	if len(p) == 0 {
		return nil, false
	}
	return p[0].parsed, true
}

// SkipTo line matching predicate
func (p Parsed) SkipTo(predicate func(string) bool) Parsed {
	for i, s := range p {
		if predicate(s.original) {
			return p[i:]
		}
	}
	return nil
}

// TakeTo removes everything after the first match of the predicate
func (p Parsed) TakeTo(predicate func(string) bool) Parsed {
	for i, s := range p {
		if predicate(s.original) {
			return p[:i]
		}
	}
	return p
}

// SkipOneLine or none if text is already empty
func (p Parsed) SkipOneLine() Parsed {
	if len(p) == 0 {
		return p
	}
	return p[1:]
}
