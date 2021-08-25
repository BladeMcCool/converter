package web

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kr/pretty"
	"io"
	"log"
	"net/http"
	"strings"
)

func runServer() {
	http.ListenAndServe(":5445", getRouter())
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/{inputFormat}/{outputFormat}/", ConvertDocument)
	return r
}

func ConvertDocument(w http.ResponseWriter, r *http.Request) {


	//mm = multipart.N

	//kak, _ := ioutil.ReadAll(r.Body)
	//log.Print("----body----\n")
	//log.Printf("%s", pretty.Formatter(string(kak)))
	//return

	//r.ParseMultipartForm(32 << 20) // limit your max input length!
	//var buf bytes.Buffer
	//// in your case file would be fileupload
	//file, header, err := r.FormFile("data")
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	//name := strings.Split(header.Filename, ".")
	//fmt.Printf("File name %s\n", name[0])
	//// Copy the file data to my buffer
	//io.Copy(&buf, file)
	//contents := buf.String()
	//fmt.Println(contents)

	//
	vars := mux.Vars(r)
	log.Printf("%# v\n", pretty.Formatter(vars))
	//
	r.ParseMultipartForm(10 * 1024 * 1024)
	log.Print("---!!!----\n")
	for k, v := range r.PostForm {
		log.Printf("%s : %s\n", k, v)
	}
	log.Print("------------\n")
	log.Printf("%# v\n", r.FormValue("inputLineSeparator"))
	log.Print("------------\n")

	var buf bytes.Buffer
	file, header, err := r.FormFile("data")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	// Copy the file data to my buffer
	io.Copy(&buf, file)
	contents := buf.String()
	fmt.Println(contents)


	////log.Printf("%s : %s\n", k, v)
	//
	//
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("horses"))

}