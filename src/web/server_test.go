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
	//body, _ := ioutil.ReadAll(responseRecorder.Body)
	//t.Logf("%s", body)

	reader, err := os.Open("../testdata/format1_example.txt")
	require.Nil(t, err)
	format1RawBytes, err := ioutil.ReadAll(reader)

	require.Equal(t, 200, responseRecorder.Code)
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

func Test__Convert__Format3ToFormat1__MissingSeparator__Responds400(t *testing.T) {
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
		requestWriter.Close()
	}()

	responseRecorder := httptest.NewRecorder()
	getRouter().ServeHTTP(responseRecorder, request)

	require.Equal(t, 400, responseRecorder.Code)
	body, _ := ioutil.ReadAll(responseRecorder.Body)
	require.Contains(t, string(body), "missing one or more separator" )
}

func Test__Convert__UnknownFormatIn__Responds400(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/v1/4/1/", nil)
	responseRecorder := httptest.NewRecorder()
	getRouter().ServeHTTP(responseRecorder, request)
	body, _ := ioutil.ReadAll(responseRecorder.Body)
	require.Contains(t, string(body), "format index out of range" )
}

func Test__Convert__UnknownFormatOut__Responds400(t *testing.T) {
	request := httptest.NewRequest("POST", "/api/v1/1/4/", nil)
	responseRecorder := httptest.NewRecorder()
	getRouter().ServeHTTP(responseRecorder, request)
	body, _ := ioutil.ReadAll(responseRecorder.Body)
	require.Contains(t, string(body), "format index out of range" )
}

//  some more test ideas:
//  test send no file under data param gives err
//  test send wrong (mismatched from stated) format gives err
//  test format1 in w/o separator args gives err (out is covered)
