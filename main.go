package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/subranag/mockaroo/spec"
)

const pathSeparator = "/"

var serverSpec *spec.ServerSpec

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if serverSpec == nil {
		http.Error(w, "server spec is nil", http.StatusInternalServerError)
		return
	}
	// now let spec just handle the request
	serverSpec.PerformMockAction(w, r)
}

func main() {

	specFilePath := flag.String("specFile", "", "path to the mock spec file <required argument>")

	flag.Parse()

	if *specFilePath == "" {
		flag.Usage()
		return
	}

	var err error
	if serverSpec, err = spec.ReadSpecFile(*specFilePath); err != nil {
		log.Fatalf("error reading spec file %v", err)
		return
	}

	if serverSpec == nil {
		log.Fatal("server spec is nil cannot start server")
		return
	}

	log.Print("starting mockaroo....")
	serverSpec.LogSpecDefinition()
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", serverSpec.Port), nil))
}
