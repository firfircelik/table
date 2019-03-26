package table

import (
	"regexp"
	"strings"

	funk "github.com/thoas/go-funk"
)

// FieldMatcher is a function type which is consumed by different table Cell finder functions.
type FieldMatcher func([]string) (string, bool)

// LineContaining returns a predicate checking whether a line contains specified tokens
func LineContaining(ss ...string) func(string) bool {
	return func(line string) bool {
		for _, s := range ss {
			if !strings.Contains(line, s) {
				return false
			}
		}
		return true
	}
}

// LineContainingSlices is the slice version of line containing
func LineContainingSlices(ls ...[]string) func(string) bool {
	m := map[string]bool{}
	for _, l := range ls {
		for _, s := range l {
			m[s] = true
		}
	}
	return LineContaining(funk.Keys(m).([]string)...)
}

// LineContainingAny behaves like line containing but any given slice is enough
func LineContainingAny(ls ...[]string) func(string) bool {
	predicates := make([]func(string) bool, len(ls))
	for i, ss := range ls {
		predicates[i] = LineContaining(ss...)
	}
	return AnyMatched(predicates...)
}

// LineContainingAnySingle behaves like LineContainingAny
// but any token is enough to match, instead of requiring a slice
func LineContainingAnySingle(ls ...string) func(string) bool {
	predicates := make([]func(string) bool, len(ls))
	for i, ss := range ls {
		predicates[i] = LineContaining(ss)
	}
	return AnyMatched(predicates...)
}

// AllAreMatched accepts lines until all predicates are matched
func AllAreMatched(pp ...func(string) bool) func(string) bool {
	index := 0
	return func(line string) bool {
		if index >= len(pp) {
			return true
		}
		if pp[index](line) {
			index++
		}
		return index >= len(pp)
	}
}

// AnyMatched accepts line any one of predicates is satisfied
func AnyMatched(pp ...func(string) bool) func(string) bool {
	return func(line string) bool {
		if len(pp) == 0 {
			return true
		}
		for _, p := range pp {
			if p(line) {
				return true
			}
		}
		return false
	}
}

// NonEmptyLine checks whether line is not empty
func NonEmptyLine() func(string) bool {
	return func(s string) bool {
		return !isWhiteSpace(s)
	}
}

// EmptyLine matches empty line
func EmptyLine() func(string) bool {
	return func(s string) bool {
		return isWhiteSpace(s)
	}
}

// LineFieldMatcher provides matchers for whole line or specific cell content
type LineFieldMatcher struct {
	Re  *regexp.Regexp
	Sep string
}

// Find returns a matching cell in a given line.
func (fm LineFieldMatcher) Find(line []string) (string, bool) {
	for _, s := range line {
		if fm.Re.MatchString(s) {
			return s, true
		}
	}
	return "", false
}

// FindLine returns a matching line.
func (fm LineFieldMatcher) FindLine(fields []string) (string, bool) {
	line := strings.Join(fields, fm.Sep)
	return line, fm.Re.MatchString(line)
}
