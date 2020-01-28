package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/subranag/mockaroo/spec"
)

const pathSeparator = "/"

var serverSpec *spec.ServerSpec

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Request Details</h1>")

	for header, values := range r.Header {
		fmt.Fprintf(w, "<h3> header %v values: %v </h3>", header, values)
	}

	pathComponents := strings.Split(r.URL.Path, pathSeparator)

	fmt.Fprintf(w, "<b> path components %v </b>", pathComponents)
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
