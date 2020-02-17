package spec

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

// LogSpecDefinition logs all the specs configured to the log use this function display
// all the mocks configured
func (spec *ServerSpec) LogSpecDefinition() {
	log.Printf("mock server will start on port : %v", spec.Port)
	for _, mock := range spec.Mocks {
		log.Println("")
		log.Printf("mock path    : %v", mock.MatchPath)
		log.Printf("mock method  : %v", mock.HTTPMethod)
		log.Printf("match params : %v", mock.MatchRequestParams)
	}
	log.Println(" ")
}

func prettyLogRequest(r *http.Request) {
	log.Println("-----------------------------------------------------")
	log.Printf("reqeust path    : %v", r.URL.Path)
	log.Printf("request from    : %v", r.RemoteAddr)
	log.Printf("request method  : %v", r.Method)
	log.Printf("HTTP version    : %v", r.Proto)
	if len(r.Header) > 0 {
		log.Println()
		log.Println("request headers:")
		log.Println("----------------")
		for header, vals := range r.Header {
			log.Printf("%v : %v", header, vals)
		}
	}
	if len(r.URL.Query()) > 0 {
		log.Println()
		log.Println("requst params:")
		log.Println("--------------")
		for param, vals := range r.URL.Query() {
			log.Printf("%v : %v", param, vals)
		}
	}
	log.Println("-----------------------------------------------------")
}

// MatchMock check to see if this mock spec is a match for the incoming request
func (mockSpec *MockSpec) MatchMock(r *http.Request) bool {
	return strings.Compare(r.URL.Path[1:], mockSpec.MatchPath) == 0
}

// PerformMockAction perfroms the HTTP mock action based on ServerSpec and writes a
// suitable response
func (spec *ServerSpec) PerformMockAction(w http.ResponseWriter, r *http.Request) {
	// log the request first
	prettyLogRequest(r)

	for _, mock := range spec.Mocks {
		if mock.MatchMock(r) {
			for header, val := range mock.Action.ResponseHeaders {
				w.Header().Set(header, val)
			}

			if mock.Action.ResponseTemplate != "" {
				tmpl, err := template.New(mock.MatchPath).Parse(mock.Action.ResponseTemplate)
				if err != nil {
					serverError(w, fmt.Sprintf("error processing template : %v error %v", mock.Action.ResponseTemplate, err))
					return
				}
				reponseModel := ResponseModel{RequestPath: r.URL.Path}
				err = tmpl.Execute(w, reponseModel)
				return
			}

			if mock.Action.ResponseFile != "" {
				responseBytes, err := ioutil.ReadFile(mock.Action.ResponseFile)
				if err != nil {
					serverError(w, fmt.Sprintf("error reading response file : %v", err))
					return
				}
				w.Write(responseBytes)
				return
			}

			if mock.Action.ResponseString != "" {
				w.Write([]byte(mock.Action.ResponseString))
				return
			}
			serverError(w, fmt.Sprintf("no response defined for path : %v", mock.MatchPath))
		}
	}
	serverError(w, fmt.Sprintf("no match found for path : %v", r.URL.Path[1:]))
	return
}

func serverError(w http.ResponseWriter, errorMessage string) {
	http.Error(w, errorMessage, http.StatusInternalServerError)
}
