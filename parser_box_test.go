package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type parserSuite struct{ suite.Suite }

func TestBox(t *testing.T) { suite.Run(t, new(parserSuite)) }

func (p *parserSuite) TestParsingBoxes() {
	input := `
- - - - - -
 | H | H
 | 1 | 2
- - - - -
h| H1|H2
1| h1|h1
___________
h| H1|H2
2| h2|h2

`
	result, err := ParseBoxes(strings.Split(input, "\n"), 3)
	require.Nil(p.T(), err)

	expectedResult := map[Key]string{
		{"H 1", "h 1"}: "H1 h1",
		{"H 2", "h 1"}: "H2 h1",
		{"H 1", "h 2"}: "H1 h2",
		{"H 2", "h 2"}: "H2 h2",
	}

	require.Equal(p.T(), expectedResult, result)
}
