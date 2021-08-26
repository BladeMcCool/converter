package web

import (
	"bytes"
	"converter/formats"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

func ConvertDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inputFormat, err := getFormat(vars, "inputFormat")
	if err != nil {
		writeErr(w, err)
		return
	}
	outputFormat, err := getFormat(vars, "outputFormat")
	if err != nil {
		writeErr(w, err)
		return
	}
	err = r.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		writeErr(w, err)
		return
	}

	file, _, err := r.FormFile("data")
	if err != nil {
		writeErr(w, err)
		return
	}
	defer file.Close()
	var buf bytes.Buffer
	io.Copy(&buf, file)

	inputImplementation := inputFormat.Implementation()
	err = inputImplementation.SetDelimiters("input", &r.PostForm)
	if err != nil {
		writeErr(w, err)
		return
	}
	segments, err := inputImplementation.Import(buf.Bytes())
	if err != nil {
		writeErr(w, err)
		return
	}

	outputImplementation := outputFormat.Implementation()
	err = outputImplementation.SetDelimiters("output", &r.PostForm)
	if err != nil {
		writeErr(w, err)
		return
	}
	exportData := outputImplementation.Export(segments)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", outputFormat.Mimetype)
	w.Write(exportData)

}

func getFormat(vars map[string]string, optionField string) (*formats.FormatSpec, error){
	inputFormatIndexStr, ok := vars[optionField]
	if !ok {
		return nil, fmt.Errorf("invalid request, missing %s param", optionField)
	}

	inputFormatIndex, err := strconv.Atoi(inputFormatIndexStr)
	if err != nil { return nil, err }
	if inputFormatIndex < 1 || inputFormatIndex > len(formats.FormatSpecs) {
		return nil, fmt.Errorf("format index out of range for %s", optionField)
	}
	return &formats.FormatSpecs[inputFormatIndex-1], nil
}

func writeErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func RunServer() {
	http.ListenAndServe(":5445", getRouter())
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/{inputFormat}/{outputFormat}/", ConvertDocument)
	return r
}
