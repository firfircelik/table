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
		{original: "location\tdelivery\tprice", parsed: []string{"location", "delivery", "price"}},
		{original: "delhi\tjanuary\t100", parsed: []string{"delhi", "january", "100"}},
		{original: "pune\tfebruary\t80", parsed: []string{"pune", "february", "80"}},
		{original: "pune\tno prices\t", parsed: []string{"pune", "no prices", ""}},
	}

	// When
	parsedHTML, err := ParseFromHTML(htmlStr)

	// Then
	require.Nil(s.T(), err)
	require.Equal(s.T(), expectedParsedHTML, parsedHTML)

	const htmlDiacriticsString = `
<html>
  <body>
    <table>
		<tr>
			<td>NĚw YoRK</td>
			<td>MaRch</td>
			<td>200</td>
		</tr>
		<tr>
			<td>ZÙriÇh</td>
			<td>April</td>
			<td>100</td>
		</tr>
		<tr>
			<td>RoMÊ</td>
			<td>JuNĚ</td>
			<td>100</td>
		</tr>
	</table>
  </body>
</html>
`

	expectedParsedHTML = Parsed{
		{original: "new york\tmarch\t200", parsed: []string{"new york", "march", "200"}},
		{original: "zurich\tapril\t100", parsed: []string{"zurich", "april", "100"}},
		{original: "rome\tjune\t100", parsed: []string{"rome", "june", "100"}},
	}

	// When
	parsedHTML, err = ParseFromHTML(htmlDiacriticsString)

	// Then
	require.Nil(s.T(), err)
	require.Equal(s.T(), expectedParsedHTML, parsedHTML)
}
