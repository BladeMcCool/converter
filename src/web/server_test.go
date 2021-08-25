package web

import (
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

func Test__Upload(t *testing.T) {
	pr, pw := io.Pipe()
	requestWriter := multipart.NewWriter(pw)
	//defer pr.Close()

	request := httptest.NewRequest("POST", "/api/v1/xxx/yyy/", pr)
	request.Header.Add("Content-Type", requestWriter.FormDataContentType())

	go func() {
		defer pw.Close()
		uploadfileWriter, err := requestWriter.CreateFormFile("data", "undddknown.txt")
		require.Nil(t, err)
		uploadfileWriter.Write([]byte("hrm"))
		requestWriter.WriteField("inputLineSeparator", "xx&xx")
		requestWriter.Close()
		//inputLineSeparator
		//inputElementSeparator
		//outputLineSeparator
		//outputElementSeparator
	}()

	responseRecorder := httptest.NewRecorder()
	//handler := ConvertDocument
	getRouter().ServeHTTP(responseRecorder, request)
	t.Logf("%# v", responseRecorder.Code)
	t.Logf("%# v", responseRecorder.Header())
	require.Equal(t, "horses", responseRecorder.Body.String())
}

