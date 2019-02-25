package table

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type parseHTMLSuite struct{ suite.Suite }

func TestParseHTML(t *testing.T) { suite.Run(t, new(parseHTMLSuite)) }

func (s *parseHTMLSuite) TestParseFromHTML() {

	// Given
	const htmlStr = `
<html>
  <body>
    <table>
      <tr>
        <th>Location</th>
        <th>Delivery</th>
        <th>Price</th>
      </tr>
      <tr>
        <td>Delhi</td>
        <td>January</td>
        <td>100</td>
      </tr>
      <tr>
        <td>Pune</td>
        <td>February</td>
        <td>80</td>
      </tr>
      <tr>
        <td>Pune</td>
        <td colspan="2">No prices</td>
      </tr>
    </table>
  </body>
</html>
`

	var expectedParsedHTML = Parsed{
		{original: "Location\tDelivery\tPrice", parsed: []string{"Location", "Delivery", "Price"}},
		{original: "Delhi\tJanuary\t100", parsed: []string{"Delhi", "January", "100"}},
		{original: "Pune\tFebruary\t80", parsed: []string{"Pune", "February", "80"}},
		{original: "Pune\tNo prices\t", parsed: []string{"Pune", "No prices", ""}},
	}

	// When
	parsedHTML, err := ParseFromHTML(htmlStr)

	// Then
	require.Nil(s.T(), err)
	require.Equal(s.T(), expectedParsedHTML, parsedHTML)
}
