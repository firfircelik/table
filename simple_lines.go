package table

import (
	"regexp"
)

var twoOrMoreWhitespaces = regexp.MustCompile(`\s\s\s*`)

// T is a not parsed table
type T []string

// SkipTo line matching predicate
func (t T) SkipTo(predicate func(string) bool) T {
	for i, s := range t {
		if predicate(s) {
			return t[i:]
		}
	}
	return nil
}

// TakeTo removes everything after the first match of the predicate
func (t T) TakeTo(predicate func(string) bool) T {
	for i, s := range t {
		if predicate(s) {
			return t[:i]
		}
	}
	return t
}

// TakeIncluding removes everything after the first match of the predicate
func (t T) TakeIncluding(predicate func(string) bool) T {
	for i, s := range t {
		if predicate(s) {
			return t[:i+1]
		}
	}
	return t
}

// SkipOneLine or none if text is already empty
func (t T) SkipOneLine() T {
	if len(t) == 0 {
		return t
	}
	return t[1:]
}

// FirstOrEmpty returns first line or empty string
func (t T) FirstOrEmpty() string {
	if len(t) == 0 {
		return ""
	}
	return t[0]
}

// Ensure validates predicates are correct at least one line in the table
// Empty table is wrong by default
func (t T) Ensure(predicates ...func(string) bool) bool {
	if len(predicates) == 0 {
		return true
	}
	for _, predicate := range predicates {
		var correct bool
		for _, line := range t {
			if predicate(line) {
				correct = true
				break
			}
		}
		if !correct {
			return false
		}
	}
	return len(t) > 0
}

// IgnoreLines removes lines from given table
func (t T) IgnoreLines(lines []string) T {
	filtered := make([]string, 0, len(t))
	for _, l := range t {
		if !containsString(lines, l) {
			filtered = append(filtered, l)
		}
	}
	return T(filtered)
}

// sliceIndex returns first index of `x` in `slice` and -1 if `x` is not present.
func sliceIndex(slice []string, x string) int {
	for i, v := range slice {
		if v == x {
			return i
		}
	}
	return -1
}

// containsString checks if array contains element
func containsString(s []string, e string) bool {
	return sliceIndex(s, e) > -1
}
