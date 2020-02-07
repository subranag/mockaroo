package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/subranag/mockaroo/spec"
)

const pathSeparator = "/"

var serverSpec *spec.ServerSpec

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if serverSpec == nil {
		http.Error(w, "server spec is nil", http.StatusInternalServerError)
		return
	}
	log.Printf("request path is : %v", r.URL.Path)

	for _, mock := range serverSpec.Mocks {
		if strings.Compare(r.URL.Path[1:], mock.MatchPath) == 0 {

			for header, val := range mock.Action.ResponseHeaders {
				w.Header().Set(header, val)
			}

			if mock.Action.ResponseTemplate != "" {
				tmpl, err := template.New(mock.MatchPath).Parse(mock.Action.ResponseTemplate)
				if err != nil {
					http.Error(w, fmt.Sprintf("error processing template : %v error %v", mock.Action.ResponseTemplate, err), http.StatusInternalServerError)
					return
				}
				reponseModel := spec.ResponseModel{RequestPath: r.URL.Path}
				err = tmpl.Execute(w, reponseModel)
				return
			}

			responseBytes, err := ioutil.ReadFile(mock.Action.ResponseFile)
			if err != nil {
				http.Error(w, fmt.Sprintf("error reading response file : %v", err), http.StatusInternalServerError)
				return
			}
			w.Write(responseBytes)
			return
		}
	}
	http.Error(w, fmt.Sprintf("no match found for path : %v", r.URL.Path[1:]), http.StatusInternalServerError)
	return
}

func logSpecDetails() {
	log.Printf("mock server will start on port : %v", serverSpec.Port)
	for _, mock := range serverSpec.Mocks {
		log.Println("")
		log.Printf("mock path    : %v", mock.MatchPath)
		log.Printf("mock method  : %v", mock.HTTPMethod)
		log.Printf("match params : %v", mock.MatchRequestParams)
	}
	log.Println(" ")
}

func main() {
	var err error
	if serverSpec, err = spec.ReadSpecFile("/var/tmp/sample_spec.yaml"); err != nil {
		log.Fatalf("error reading spec file %v", err)
		return
	}

	if serverSpec == nil {
		log.Fatal("server spec is nil cannot start server")
		return
	}

	logSpecDetails()

	log.Print("starting mockaroo....")
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
