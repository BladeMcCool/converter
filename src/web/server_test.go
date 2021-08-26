package web

import (
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test__Convert__Format3ToFormat1__OutputCorrect(t *testing.T) {
	pr, pw := io.Pipe()
	requestWriter := multipart.NewWriter(pw)

	request := httptest.NewRequest("POST", "/api/v1/3/1/", pr)
	request.Header.Add("Content-Type", requestWriter.FormDataContentType())

	go func() {
		defer pw.Close()
		uploadfileWriter, err := requestWriter.CreateFormFile("data", "allrecords.xml")
		require.Nil(t, err)
		reader, err := os.Open("../testdata/format3_example.xml")
		require.Nil(t, err)
		fixtureRawBytes, err := ioutil.ReadAll(reader)

		uploadfileWriter.Write(fixtureRawBytes)
		requestWriter.WriteField("outputLineSeparator", "~")
		requestWriter.WriteField("outputElementSeparator", "*")
		requestWriter.Close()
		//inputLineSeparator
		//inputElementSeparator
	}()

	responseRecorder := httptest.NewRecorder()
	getRouter().ServeHTTP(responseRecorder, request)
	t.Logf("%# v", responseRecorder.Code)
	t.Logf("%# v", responseRecorder.Header())

	reader, err := os.Open("../testdata/format1_example.txt")
	require.Nil(t, err)
	format1RawBytes, err := ioutil.ReadAll(reader)

	require.Equal(t, "text/plain", responseRecorder.Header().Get("Content-Type"))

	//we said the samples of the formats were equivalent
	//we noted whitespace differences and thus decided to trim during import
	//we arent tracking where whitespace was trimmed or which fields of which formats would need it added back
	//so output here from import->export will have those spaces stripped
	//if we need to account for different spacing in certain fields for certain formats and not others
	//then that logic must be described first.
	//
	//in the meantime, to compare the output vs the fixture i'm going to strip all spaces and count it as correct if they match.
	format1RawBytesNoWhitespace := strings.ReplaceAll(string(format1RawBytes), " ", "")
	require.Equal(t, format1RawBytesNoWhitespace, responseRecorder.Body.String())
}

//  test send no file under data param gives err
//  test send wrong format gives err
//  test unknown format in gives err
//  test unknown format out gives err
//  test format1 in w/o separator args gives err
//  test format1 out w/o separator args gives err
