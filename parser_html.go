package table

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/agflow/agstring"
	"github.com/pkg/errors"
)

func getHTMLOptions(options ...*Options) *Options {
	opts := getOptions(options...)
	if agstring.IsEmpty(opts.RowSelector) {
		opts.RowSelector = "table tr"
	}
	return opts
}

// ParseFromHTML parses HTML tables into string with a variety of customizations
// default selector: "table tr"
func ParseFromHTML(s string, options ...*Options) (Parsed, error) {
	opts := getHTMLOptions(options...)

	s = agstring.NormalizeDiacritics(s)
	var p Parsed
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return p, errors.Wrap(err, "can't parse html")
	}
	doc.Find(opts.RowSelector).Each(func(i int, s *goquery.Selection) {
		if err != nil {
			return
		}
		var line []string
		s.Find("td,th").Each(func(i int, s *goquery.Selection) {
			if err != nil {
				return
			}
			colspan, err2 := getColspan(s, opts)
			if err2 != nil {
				err = err2
			}
			line = append(line, strings.TrimSpace(s.Text()))
			for i := 1; i < colspan; i++ {
				line = append(line, "")
			}
		})
		p = append(p, parsedLine{
			original: strings.Join(line, "\t"),
			parsed:   line,
		})
	})
	return p, err
}

func getColspan(s *goquery.Selection, opts *Options) (int, error) {
	colspan, err := extractColspan(s)
	if err != nil {
		return 0, err
	}
	if opts.IgnoreColspan {
		colspan = 1
	}
	return colspan, nil
}

func extractColspan(s *goquery.Selection) (int, error) {
	val, ok := s.Attr("colspan")
	if !ok {
		return 1, nil
	}
	n, err := strconv.Atoi(val)
	return n, errors.Wrap(err, "can't parse colspan")
}
