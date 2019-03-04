package agstring

import (
	"strings"
)

// Normalize first lowercase string and then trim it
func Normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// NormalizeDiacritics remove diacritics from a normalized string
func NormalizeDiacritics(s string) string {
	return RemoveDiacritics(Normalize(s))
}

// NormalizeDiacriticsAndNonAlnum first remove diacritics from a normalized string then remove
// all non alphanumeric characters including whitespaces in the middle
func NormalizeDiacriticsAndNonAlnum(s string) string {
	return RemoveNonAlnum(NormalizeDiacritics(s))
}
