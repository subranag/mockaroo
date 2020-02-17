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
	log.Printf("HTTP version    : %v", r.Proto)
	log.Println()
	log.Println("request headers:")
	log.Println("----------------")
	for header, vals := range r.Header {
		log.Printf("%v : %v", header, vals)
	}
	log.Println()
	log.Println("requst params:")
	log.Println("--------------")
	for param, vals := range r.URL.Query() {
		log.Printf("%v : %v", param, vals)
	}
	log.Println("-----------------------------------------------------")
}

// PerformMockAction perfroms the HTTP mock action based on ServerSpec and writes a
// suitable response
func (spec *ServerSpec) PerformMockAction(w http.ResponseWriter, r *http.Request) {
	// log the request first
	prettyLogRequest(r)

	for _, mock := range spec.Mocks {
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
				reponseModel := ResponseModel{RequestPath: r.URL.Path}
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
