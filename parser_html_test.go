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
	var htmlStr = `
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

func (s *parseHTMLSuite) TestParseFromHTMLRowSelectorOption() {

	// Given
	var htmlStr = `
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
  <div> Test </div>
  <table class="pretty">
    <tr>
      <th>Pretty Location</th>
      <th>Pretty Delivery</th>
      <th>Pretty Price</th>
    </tr>
    <tr>
      <td>Orlando</td>
      <td>April</td>
      <td>550</td>
    </tr>
  </table>
</body>
</html>
`
	expectedParsedHTML := Parsed{
		{
			original: "pretty location\tpretty delivery\tpretty price",
			parsed:   []string{"pretty location", "pretty delivery", "pretty price"},
		}, {
			original: "orlando\tapril\t550",
			parsed:   []string{"orlando", "april", "550"},
		},
	}

	// When
	parsedHTML, err := ParseFromHTML(htmlStr, &Options{RowSelector: "table.pretty tr"})

	// Then
	require.Nil(s.T(), err)
	require.Equal(s.T(), expectedParsedHTML, parsedHTML)

	// Given
	htmlStr = `
<html>
 <body>
   <table>
     <tr>
       <td>Delhi</td>
       <td>January</td>
       <td>100</td>
     </tr>
   </table>
   <div> Test </div>
   <table>
     <tr>
       <td>Orlando</td>
       <td>April</td>
       <td>550</td>
     </tr>
   </table>
   <div> Test </div>
   <table id='idtable'>
     <tr>
       <td>New York</td>
       <td>December</td>
       <td>1056</td>
     </tr>
   </table>
 </body>
</html>
`
	expectedParsedHTML = Parsed{
		{
			original: "new york\tdecember\t1056",
			parsed:   []string{"new york", "december", "1056"},
		},
	}

	// When
	parsedHTML, err = ParseFromHTML(htmlStr, &Options{RowSelector: "table:nth-of-type(3) tr"})

	// Then
	require.Nil(s.T(), err)
	require.Equal(s.T(), expectedParsedHTML, parsedHTML)
}

func (s *parseHTMLSuite) TestParseFromHTMLColspanOption() {

	// Given
	var tableWithColspan = `
<html>
  <body>
    <table>
		<tr>
			<td colspan="2">NĚw YoRK</td>
			<td>MaRch</td>
			<td>200</td>
		</tr>
		<tr>
			<td colspan="2">ZÙriÇh</td>
			<td>April</td>
			<td>100</td>
		</tr>
		<tr>
			<td colspan="2"">RoMÊ</td>
			<td>JuNĚ</td>
			<td>100</td>
		</tr>
	</table>
  </body>
</html>
`

	expectedParsedHTML := Parsed{
		{original: "new york\t\tmarch\t200", parsed: []string{"new york", "", "march", "200"}},
		{original: "zurich\t\tapril\t100", parsed: []string{"zurich", "", "april", "100"}},
		{original: "rome\t\tjune\t100", parsed: []string{"rome", "", "june", "100"}},
	}

	// When
	parsedHTML, err := ParseFromHTML(tableWithColspan)

	// Then
	require.Nil(s.T(), err)
	require.Equal(s.T(), expectedParsedHTML, parsedHTML)

	// To Ignore colspan
	expectedParsedHTML = Parsed{
		{original: "new york\tmarch\t200", parsed: []string{"new york", "march", "200"}},
		{original: "zurich\tapril\t100", parsed: []string{"zurich", "april", "100"}},
		{original: "rome\tjune\t100", parsed: []string{"rome", "june", "100"}},
	}

	// When
	parsedHTML, err = ParseFromHTML(tableWithColspan, &Options{IgnoreColspan: true})

	// Then
	require.Nil(s.T(), err)
	require.Equal(s.T(), expectedParsedHTML, parsedHTML)

	// Given
	tableWithColspan = `
<html>
  <body>
    <table>
		<tr>
			<td colspan="2">NĚw YoRK</td>
			<td>MaRch</td>
			<td>200</td>
		</tr>
		<tr>
			<td>ZÙriÇh</td>
			<td>Three</td>
			<td>April</td>
			<td>100</td>
		</tr>
		<tr>
			<td colspan="2"">RoMÊ</td>
			<td>JuNĚ</td>
			<td>100</td>
		</tr>
	</table>
  </body>
</html>
`

	expectedParsedHTML = Parsed{
		{original: "new york\t\tmarch\t200", parsed: []string{"new york", "", "march", "200"}},
		{original: "zurich\tthree\tapril\t100", parsed: []string{"zurich", "three", "april", "100"}},
		{original: "rome\t\tjune\t100", parsed: []string{"rome", "", "june", "100"}},
	}

	// When
	parsedHTML, err = ParseFromHTML(tableWithColspan)

	// Then
	require.Nil(s.T(), err)
	require.Equal(s.T(), expectedParsedHTML, parsedHTML)

	expectedParsedHTML = Parsed{
		{original: "new york\tmarch\t200", parsed: []string{"new york", "march", "200"}},
		{original: "zurich\tthree\tapril\t100", parsed: []string{"zurich", "three", "april", "100"}},
		{original: "rome\tjune\t100", parsed: []string{"rome", "june", "100"}},
	}

	// When
	parsedHTML, err = ParseFromHTML(tableWithColspan, &Options{IgnoreColspan: true})

	// Then
	require.Nil(s.T(), err)
	require.Equal(s.T(), expectedParsedHTML, parsedHTML)
}
