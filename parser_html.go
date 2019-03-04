package table

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/agflow/agstring"
	"github.com/pkg/errors"
)

// ParseFromHTML table encoded inside string
func ParseFromHTML(s string) (Parsed, error) {
	p := Parsed{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return p, errors.Wrap(err, "can't parse html")
	}
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		if err != nil {
			return
		}
		originalLine := []string{}
		line := []string{}
		s.Find("td,th").Each(func(i int, s *goquery.Selection) {
			if err != nil {
				return
			}
			colspan, err2 := extractColspan(s)
			if err2 != nil {
				err = err2
			}
			originalLine = append(originalLine, strings.TrimSpace(s.Text()))
			line = append(line, string(agstring.NormalizeDiacritics(s.Text())))
			for i := 1; i < colspan; i++ {
				line = append(line, "")
			}
		})
		p = append(p, parsedLine{
			original: strings.Join(originalLine, "\t"),
			parsed:   line,
		})
	})
	return p, err
}

func extractColspan(s *goquery.Selection) (int, error) {
	val, ok := s.Attr("colspan")
	if !ok {
		return 1, nil
	}
	n, err := strconv.Atoi(val)
	return n, errors.Wrap(err, "can't parse colspan")
}
