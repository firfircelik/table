package agstring

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mozillazg/go-unidecode"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

// ReplaceMultispace replaces multiple spaces with one space and
// also trims space from both ends
func ReplaceMultispace(s string) string {
	stripper := regexp.MustCompile(`\s{2,}`)
	return strings.TrimSpace(stripper.ReplaceAllString(s, " "))
}

// Nth returns nth element of given slice or empty string if out of limits
func Nth(ls []string, n int) string {
	if len(ls) == 0 || n < 0 || n >= len(ls) {
		return ""
	}
	return ls[n]
}

// First returns the first element of given list or empty string when the list is empty.
func First(ls ...string) string { return Nth(ls, 0) }

// Last returns the last element of given list or empty string when the list is empty.
func Last(ls ...string) string { return Nth(ls, len(ls)-1) }

// TrimSuffixes returns s without any of the provided trailing suffixes strings.
func TrimSuffixes(s string, suffixes ...string) string {
	s = strings.TrimSpace(s)
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return strings.TrimSpace(strings.TrimSuffix(s, suffix))
		}
	}
	return s
}

// TrimAllSuffixes returns a string without any of the provided trailing suffixes or spaces.
// See test for examples.
func TrimAllSuffixes(s string, suffixes ...string) string {
	if len(suffixes) == 0 || s == "" {
		return strings.TrimSpace(s)
	}

	reSufs := make([]*regexp.Regexp, 0)
	for _, suffix := range suffixes {
		if suffix == "" {
			continue
		}

		reE := fmt.Sprintf(`^(?P<rest>.*)%s\s*$`, regexp.QuoteMeta(suffix))
		reSufs = append(reSufs, regexp.MustCompile(reE))
	}

	trimAgain := true
	for trimAgain {
		trimAgain = false
		for _, reSuf := range reSufs {
			if matches, ok := RegexpGroups(reSuf, strings.TrimSpace(s)); ok {
				s = matches["rest"]
				trimAgain = true
			}
		}
	}
	return strings.TrimSpace(s)
}

// TrimAllPrefixes returns a string without any of the provided leading prefixes or spaces.
// See test for examples.
func TrimAllPrefixes(s string, prefixes ...string) string {
	if len(prefixes) == 0 || s == "" {
		return strings.TrimSpace(s)
	}

	rePres := make([]*regexp.Regexp, 0)
	for _, prefix := range prefixes {
		if prefix == "" {
			continue
		}

		reE := fmt.Sprintf(`^\s*%s(?P<rest>.*)`, regexp.QuoteMeta(prefix))
		rePres = append(rePres, regexp.MustCompile(reE))
	}

	trimAgain := true
	for trimAgain {
		trimAgain = false
		for _, rePre := range rePres {
			if matches, ok := RegexpGroups(rePre, strings.TrimSpace(s)); ok {
				s = matches["rest"]
				trimAgain = true
			}
		}
	}
	return strings.TrimSpace(s)
}

// TrimPrefixesAndSpace returns a string without any of the provided leading prefixes at word
// boundaries or spaces. See test for examples.
func TrimPrefixesAndSpace(s string, prefixes ...string) string {
	if prefixes == nil || s == "" {
		return s
	}

	rePres := make([]*regexp.Regexp, 0)
	for _, prefix := range prefixes {
		if prefix == "" {
			continue
		}

		reE := fmt.Sprintf("^\\s*%s\\b(?P<rest>.*)", regexp.QuoteMeta(prefix))
		rePres = append(rePres, regexp.MustCompile(reE))
	}

	trimAgain := true
	for trimAgain {
		trimAgain = false
		for _, rePre := range rePres {
			if matches, ok := RegexpGroups(rePre, strings.TrimSpace(s)); ok {
				s = matches["rest"]
				trimAgain = true
			}
		}
	}
	return strings.TrimSpace(s)
}

var nonAlphanumRegexp = regexp.MustCompile("[^[:alnum:]]")

// RemoveNonAlnum removes non-alphanumeric characters from string
func RemoveNonAlnum(s string) string {
	return nonAlphanumRegexp.ReplaceAllLiteralString(s, "")
}

// ContainsAll checks if given slice contains all searched strings
func ContainsAll(holder []string, searched ...string) bool {
	for _, s := range searched {
		if !funk.ContainsString(holder, s) {
			return false
		}
	}
	return true
}

// StringContainsAll checks if given string contains all searched strings
func StringContainsAll(holder string, searched ...string) bool {
	for _, s := range searched {
		if !strings.Contains(holder, s) {
			return false
		}
	}
	return true
}

// ContainsAny checks if source slice contains any of given strings
func ContainsAny(src []string, qs ...string) bool {
	for _, q := range qs {
		if funk.ContainsString(src, q) {
			return true
		}
	}
	return false
}

// StringContainsAny is similar to ContainsAny but source is a string
func StringContainsAny(s string, ls ...string) bool {
	for _, e := range ls {
		if strings.Contains(s, e) {
			return true
		}
	}
	return false
}

// RemoveDiacritics removes diacritics from a string. If non-alphanumeric character is
// encountered diacritics are removed from it. If removing diacritics is not possible, character
// is removed.
func RemoveDiacritics(s string) string { return unidecode.Unidecode(s) }

// EmptyIf returns empty string if given string equals to one
// of the strings in empty list. Otherwise, given string is returned as it is.
func EmptyIf(s string, emptyList ...string) string {
	return ConvertIf(s, "", emptyList...)
}

// ConvertIf returns converted string if given string is one of the strings in the list
func ConvertIf(val, converted string, list ...string) string {
	for _, t := range list {
		if val == t {
			return converted
		}
	}
	return val
}

// ValueIfExists returns value from map for given key if exists, else returns the given key
func ValueIfExists(k string, m map[string]string) string {
	v, ok := m[k]
	if ok {
		return v
	}
	return k
}

// ReplaceWholeWord replaces old into new only if old occurs as a whole word.
func ReplaceWholeWord(s, old, replacement string) string {
	s = " " + s + " "
	old = " " + old + " "
	replacement = " " + replacement + " "
	s = strings.Replace(s, old, replacement, -1)
	return s[1 : len(s)-1]
}

// StringIterator provides a generator of names / strings
type StringIterator interface {
	Get() string
	HasNext() bool
}

// TrimSpace trims spaces in the given slice
func TrimSpace(ls ...string) []string { return Map(ls, strings.TrimSpace) }

// ToLower makes lowercase strings in the given slice
func ToLower(ls ...string) []string { return Map(ls, strings.ToLower) }

// Title ensures title formatting for given string
func Title(s string) string { return strings.Title(strings.ToLower(s)) }

// HasSuffix checks string has any one of given suffixes
func HasSuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

// HasPrefix checks string has any one of given prefixes
func HasPrefix(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// NonEmpty filters nonempty strings from given slice
func NonEmpty(ls ...string) []string {
	return nonempty(ls, UnitTransform)
}

// NonEmptyIfTrimmed filters nonempty string only if
// contains some data when whitespace is removed
func NonEmptyIfTrimmed(ls ...string) []string {
	return nonempty(ls, strings.TrimSpace)
}

// Transform modifies a string to another format
type Transform func(string) string

// UnitTransform doesn't modify its argument
func UnitTransform(s string) string { return s }

func nonempty(ls []string, t Transform) []string {
	var nonempty []string
	for _, s := range ls {
		if t(s) != "" {
			nonempty = append(nonempty, s)
		}
	}
	return nonempty
}

// IsEmpty checks if slice contains only empty strings
func IsEmpty(ls ...string) bool { return len(NonEmpty(ls...)) == 0 }

// RemoveAllDiacritics removes diacritics from all strings in slice
func RemoveAllDiacritics(ls ...string) []string { return Map(ls, RemoveDiacritics) }

// SafeAtoi converts string, including empty string, to int
func SafeAtoi(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(s)
	return n, errors.Wrap(err, "can't convert to int")
}

// RegexpGroups checks if regex matches to given string
// If so, returns named groups with matches in a map
func RegexpGroups(exp *regexp.Regexp, input string) (map[string]string, bool) {
	if !exp.MatchString(input) {
		return nil, false
	}
	match := exp.FindStringSubmatch(input)
	result := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result, true
}

// TakeTo truncates each string in the input slice up to `n` characters.
func TakeTo(ls []string, n int) []string {
	out := make([]string, 0, len(ls))
	for _, s := range ls {
		rs := []rune(s)
		o := string(rs[:min(len(rs), n)])
		out = append(out, o)
	}
	return out
}

// TakeFrom removes the first `n` characters from each string in the input slice
func TakeFrom(ls []string, n int) []string {
	out := make([]string, 0, len(ls))
	for _, s := range ls {
		rs := []rune(s)
		o := string(rs[min(len(rs), n):])
		out = append(out, o)
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ReplaceDayOrdinal replaces day ordinals (`st`, `nd`, `rd`, `th`)
// Default replaces with empty string.
func ReplaceDayOrdinal(s string, replacements ...string) string {
	var rep string
	if len(replacements) > 0 {
		rep = replacements[0]
	}
	ordinal := strings.NewReplacer("st", rep, "nd", rep, "th", rep, "rd", rep)
	return ordinal.Replace(s)
}

// ReplaceNewline replaces the newline character `\n`
// Default replaces with empty string.
func ReplaceNewline(s string, replacements ...string) string {
	var rep string
	if len(replacements) > 0 {
		rep = replacements[0]
	}
	return strings.Replace(s, "\n", rep, -1)
}

// Map runs given modifiers for each item in slice and returns a new slice
func Map(ls []string, funcs ...func(string) string) []string {
	out := make([]string, len(ls))
	for i, s := range ls {
		tmp := s
		for _, f := range funcs {
			tmp = f(tmp)
		}
		out[i] = tmp
	}
	return out
}

func Concat(slices ...[]string) []string {
	var concatSlice []string
	for _, slice := range slices {
		concatSlice = append(concatSlice, slice...)
	}
	return concatSlice
}
