package formats

import (
	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

//note:
// given differences in leading and trailing spaces between some corresponding values in the 3 formats
// and the fact that they are all deemed equivalent we have to surmise that leading and trailing spaces
// are not important and should be stripped from internal format.
// also meaning that comparing export to the fixtures we have to ignore any difference in leading/trailing spaces on the values.
var importedSegmentsFixture = []Segment{{
	kind:   "ProductID",
	values: []string{
		"4", "8", "15", "16", "23",
	},
},{
	kind:   "AddressID",
	values: []string{
		"42", "108", "3", "14",
	},
},{
	kind:   "ContactID",
	values: []string{
		"59", "26",
	},
}}

var exportMatchJSONAndXMLFixture = []Segment{{
	kind:   "ProductID",
	values: []string{
		"4", "8", "15", "16", "23",
	},
},{
	kind:   "AddressID",
	values: []string{
		"42", "108", "3", " 14",
	},
},{
	kind:   "ContactID",
	values: []string{
		"59", "26",
	},
}}

var exportMatchTxtFixture = []Segment{{
	kind:   "ProductID",
	values: []string{
		" 4", " 8", " 15", " 16", "23",
	},
},{
	kind:   "AddressID",
	values: []string{
		" 42", " 108", "3", " 14 ",
	},
},{
	kind:   "ContactID",
	values: []string{
		" 59", " 26",
	},
}}

func TestFormat1__Import__MatchesExpected(t *testing.T)  {
	reader, err := os.Open("../testdata/format1_example.txt")
	require.Nil(t, err)
	fixtureRawBytes, err := ioutil.ReadAll(reader)
	formatter := Format1Formatter{}
	formatter.SetDelimiters("~", "*")
	importedData, err := formatter.Import(fixtureRawBytes)
	t.Logf("%# v\n", pretty.Formatter(importedData))

	require.Nil(t, err)
	require.Equal(t, importedSegmentsFixture, importedData)
}

func TestFormat1__Import__BadData__ReportsInvalidInput(t *testing.T) {
	reader, err := os.Open("../testdata/format1_baddata.txt")
	require.Nil(t, err)
	fixtureRawBytes, err := ioutil.ReadAll(reader)
	formatter := Format1Formatter{}
	formatter.SetDelimiters("~", "*")
	_, err = formatter.Import(fixtureRawBytes)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid input")
}

func TestFormat1__Export__MatchesExpected(t *testing.T)  {
	formatter := Format1Formatter{}
	formatter.SetDelimiters("~", "*")
	exportedData := formatter.Export(exportMatchTxtFixture)

	reader, err := os.Open("../testdata/format1_example.txt")
	require.Nil(t, err)
	fixtureRawBytes, _ := ioutil.ReadAll(reader)

	require.Equal(t, string(fixtureRawBytes), string(exportedData))
}

