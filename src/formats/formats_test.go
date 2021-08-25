package formats

import (
	"encoding/json"
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
		" 4", " 8", " 15", " 16", "23",
	},
},{
	kind:   "AddressID",
	values: []string{
		" 42", " 108", "3", " 14",
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
	reader, err := os.Open("../testdata/format1_example2.txt")
	require.Nil(t, err)
	fixtureRawBytes, err := ioutil.ReadAll(reader)
	formatter := Format1ImporterExporter{}
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
	formatter := Format1ImporterExporter{}
	formatter.SetDelimiters("~", "*")
	_, err = formatter.Import(fixtureRawBytes)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid input")
}

func TestFormat1__Export__MatchesExpected(t *testing.T)  {
	formatter := Format1ImporterExporter{}
	formatter.SetDelimiters("~", "*")
	exportedData := formatter.Export(exportMatchTxtFixture)

	reader, err := os.Open("../testdata/format1_example.txt")
	require.Nil(t, err)
	fixtureRawBytes, _ := ioutil.ReadAll(reader)

	require.Equal(t, string(fixtureRawBytes), string(exportedData))
}

func TestFormat2__Import__MatchesExpected(t *testing.T)  {
	reader, err := os.Open("../testdata/format2_example.json")
	require.Nil(t, err)
	fixtureRawBytes, err := ioutil.ReadAll(reader)
	formatter := Format2ImporterExporter{}
	importedData, err := formatter.Import(fixtureRawBytes)
	t.Logf("%# v\n", pretty.Formatter(importedData))

	require.Nil(t, err)
	require.Equal(t, importedSegmentsFixture, importedData)
}

func TestFormat2__Import__BadData__ReportsInvalidInput(t *testing.T) {
	reader, err := os.Open("../testdata/format2_baddata.json")
	require.Nil(t, err)
	fixtureRawBytes, err := ioutil.ReadAll(reader)
	formatter := Format2ImporterExporter{}
	_, err = formatter.Import(fixtureRawBytes)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid input")

	reader, err = os.Open("../testdata/format2_baddata2.json")
	require.Nil(t, err)
	fixtureRawBytes, err = ioutil.ReadAll(reader)
	_, err = formatter.Import(fixtureRawBytes)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unexpected end of JSON input")

	reader, err = os.Open("../testdata/format2_baddata3.json")
	require.Nil(t, err)
	fixtureRawBytes, err = ioutil.ReadAll(reader)
	_, err = formatter.Import(fixtureRawBytes)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "json: cannot unmarshal")
}


func TestFormat2__Export__MatchesExpected(t *testing.T)  {
	formatter := Format2ImporterExporter{}
	exportedData := formatter.Export(exportMatchJSONAndXMLFixture)

	reader, err := os.Open("../testdata/format2_example.json")
	require.Nil(t, err)
	fixtureRawBytes, _ := ioutil.ReadAll(reader)

	//because JSON object keys are not ordered, order here shouldnt matter and to be able to compare
	// our function output vs the expected fixture we'll need to put the fixture in the lexical order
	// that the json marshaller outputs map keys in (otherwise, we have to go deeper)
	fixtureImported := jsonFormat{}
	err = json.Unmarshal(fixtureRawBytes, &fixtureImported)
	fixtureRemarshalled, err := json.Marshal(fixtureImported)

	require.Equal(t, string(fixtureRemarshalled), string(exportedData))
}
