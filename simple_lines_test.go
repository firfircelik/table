package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type tableSuite struct{ suite.Suite }

var table = T([]string{"abc", "def", "ghi"})
var tableEmpty = T([]string{})

func TestTable(t *testing.T) { suite.Run(t, new(tableSuite)) }

func (s *tableSuite) TestSkipToWhenResult() {
	require.Equal(s.T(), T([]string{"def", "ghi"}), table.SkipTo(LineContaining("d")))
}

func (s *tableSuite) TestSkipToNoResult() {
	require.Equal(s.T(), T(nil), table.SkipTo(LineContaining("j")))
}

func (s *tableSuite) TestTakeToWhenThereIsMatch() {
	require.Equal(s.T(), T([]string{"abc"}), table.TakeTo(LineContaining("d")))
}

func (s *tableSuite) TestTakeToWhenNoMatch() {
	require.Equal(s.T(), T([]string{"abc", "def", "ghi"}), table.TakeTo(LineContaining("j")))
}

func (s *tableSuite) TestIgnoreLines() {
	require.Equal(s.T(), T([]string{"abc"}), table.IgnoreLines([]string{"def", "ghi"}))
}

func (s *tableSuite) TestIgnoreLinesNoMatch() {
	require.Equal(s.T(), table, table.IgnoreLines([]string{"xyz"}))
}

func (s *tableSuite) TestIgnoreLinesEmptyIgnore() {
	require.Equal(s.T(), table, table.IgnoreLines([]string{}))
}

func (s *tableSuite) TestIgnoreLinesEmptyTable() {
	require.Equal(s.T(), tableEmpty, tableEmpty.IgnoreLines([]string{"xyz"}))
}

func (s *tableSuite) TestIgnoreLinesEmptyTableAndIgnore() {
	require.Equal(s.T(), tableEmpty, tableEmpty.IgnoreLines([]string{}))
}

func (s *tableSuite) TestFields() {
	testcases := []struct {
		input  string
		result []field
	}{
		{
			input: "a b",
			result: []field{
				{"a b", 0}},
		}, {
			input: "a  b",
			result: []field{
				{"a", 0},
				{"b", 3}},
		}, {
			input: "a  b ",
			result: []field{
				{"a", 0},
				{"b ", 3}},
		}, {
			input: "a  b  ",
			result: []field{
				{"a", 0},
				{"b", 3}},
		}, {
			input: "  a  b",
			result: []field{
				{"a", 2},
				{"b", 5}},
		}, {
			input: "a  b  c",
			result: []field{
				{"a", 0},
				{"b", 3},
				{"c", 6}},
		}, {
			input: "a  b c  d",
			result: []field{
				{"a", 0},
				{"b c", 3},
				{"d", 8}},
		},
	}
	for _, t := range testcases {
		require.Equal(s.T(), t.result, extractFields(t.input))
	}
}

func (s *tableSuite) TestColumnOnRealData() {
	// TODO: This test shows that parser does not deal well with runes which are longer than one
	// byte. parser should be reworked to correctly with UTF.
	lines := `
Grain          Bids         Change (¢/bu)           Basis            Change
                          NOT ON THE RIVER
SRW Wheat  4.6450-4.6550       DN 0.75            7H to 8H            UNCH
Corn       3.6575-3.6775       DN 1.5             7H to 9H            UNCH
Soybeans      8.5175           DN 0.75              -21H              UNCH

                            ON THE RIVER
SRW Wheat  4.6450-4.7250       DN 0.75            7H to 15H           UNCH
Corn       3.6775-3.7075       DN 1.5             9H to 12H           UNCH
Soybeans   8.5175-8.5575       DN 0.75          -21H to -17H          UNCH`
	result, err := columns(strings.Split(lines, "\n"), 5)
	require.Nil(s.T(), err)
	require.Equal(s.T(), []column{{0, 9}, {11, 24}, {26, 42}, {48, 60}, {70, 76}}, result)
}

func (s *tableSuite) TestColumnOnSingleColumn() {
	lines := []string{"a", "abc", "defg"}
	result, err := columns(lines, 1)
	require.Nil(s.T(), err)
	require.Equal(s.T(), []column{{0, 4}}, result)
}

func (s *tableSuite) TestColumnOnThreeColumn() {
	lines := `
a    a     b
b   abc    d
g   ghjkz  e`
	result, err := columns(strings.Split(lines, "\n"), 3)
	require.Nil(s.T(), err)
	require.Equal(s.T(), []column{{0, 1}, {4, 9}, {11, 12}}, result)
}

func (s *tableSuite) TestColumnWhenOverlap() {
	lines := `
aaaaaaa    a     b
b   abc    d
g   ghjkz  e`
	_, err := columns(strings.Split(lines, "\n"), 3)
	require.NotNil(s.T(), err)
}

func (s *tableSuite) TestColumnWhenAlmostOverlap() {
	lines := `
aaaa  a
b    abc
g    ghjkz`
	result, err := columns(strings.Split(lines, "\n"), 2)
	require.Nil(s.T(), err)
	require.Equal(s.T(), []column{{0, 4}, {5, 10}}, result)
}

// nolint: dupl
func (s *tableSuite) TestParseAligned() {
	lines := `
aaa  bb   ccc
aa    bb  cc`
	result, err := ParseAligned(strings.Split(lines, "\n"), 3)
	require.Nil(s.T(), err)
	require.Equal(s.T(), [][]string{
		{"", "", ""}, // first line is empty
		{"aaa ", " bb  ", " ccc"},
		{"aa  ", "  bb ", " cc"},
	}, result.Lines())
}

func (s *tableSuite) TestParseAlignedExtendsColumnIfPossible() {
	lines := `
aax bb
aa  bb`
	result, err := ParseAligned(strings.Split(lines, "\n"), 2)
	require.Nil(s.T(), err)
	require.Equal(s.T(), [][]string{
		{"", ""},       // first line is empty
		{"aax", " bb"}, // it was possible to extend the first column
		{"aa ", " bb"},
	}, result.Lines())
}

func (s *tableSuite) TestParseAlignedGuessColumnsFromLinesWithWrongLength() {
	lines := `
aax bb  ccc  ddd
aa  bb  cc  dd
aaaaa   cc ddddd`
	result, err := ParseAligned(strings.Split(lines, "\n"), 4)
	require.Nil(s.T(), err)
	require.Equal(s.T(), [][]string{
		{"", "", "", ""},
		{"aax", " bb ", " ccc", " ddd"}, // it was possible to extend the first column
		{"aa ", " bb ", " cc ", "dd"},
		{"aaa", "aa  ", " cc ", "dddd"}, // extending aa to aaaaa would overlap the second column
	}, result.Lines())
}

// nolint: dupl
func (s *tableSuite) TestParseAlignedColumnsAreExtendedTowardsLineEnd() {
	lines := `
11  22  33
11 222 3333`
	result, err := ParseAligned(strings.Split(lines, "\n"), 3)
	require.Nil(s.T(), err)
	require.Equal(s.T(), [][]string{
		{"", "", ""},
		{"11 ", " 22 ", " 33"},
		{"11 ", "222 ", "3333"},
	}, result.Lines())
}
