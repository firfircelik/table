package agstring

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

// TakeBetween gets the string from position `from` up to `to`
// from each string in the input slice
func TakeBetween(ls []string, from, to int) []string {
	return TakeFrom(TakeTo(ls, to), from)
}

// TakeAround cuts strings up to 'to' and from 'from' and returns the combination
func TakeAround(ls []string, to, from int) []string {
	out := make([]string, 0, len(ls))
	for _, s := range ls {
		rs := []rune(s)
		l := len(rs)
		o := string(rs[:min(l, to)]) + string(rs[min(l, from):])
		out = append(out, o)
	}
	return out
}
