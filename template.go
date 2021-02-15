package mockaroo

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

const nicePrime = 2011

var stableRandom = rand.New(rand.NewSource(nicePrime))

type TemplateContext struct {
	//Method is the HTTP method for the request (GET, POST, PUT, etc.)
	Method *string

	//Protocol HTTP 1.1/2
	Protocol *string

	//Host to which the request came
	Host *string

	//RemoteAddr from which the request came
	RemoteAddr *string

	//Headers are key value pairs from request headers
	Headers http.Header

	//FormData is the data from request form and urls
	FormData url.Values

	//PathVars is the path variables captured as a part of the path
	PathVars map[string]string

	// bytes of pseudo random UUID
	uuid []byte
}

//NewUUID generates a new pseudo random UUID seeded by the same constance value
func (tc *TemplateContext) NewUUID() string {
	// generate stable pseudo GUUID from random
	stableRandom.Read(tc.uuid)
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		tc.uuid[:4],
		tc.uuid[4:6],
		tc.uuid[6:8],
		tc.uuid[8:10],
		tc.uuid[10:])
}

func (tc *TemplateContext) PathVariable(key string) string {
	if tc.PathVars != nil && tc.PathVars[key] != "" {
		return tc.PathVars[key]
	}
	return ""
}

//RandomInt generates a new pseudo random int between min and max seeded by the same constance value
func (tc *TemplateContext) RandomInt(min, max int) int {
	return min + stableRandom.Intn(max-min)
}

//RandomFloat generates a new pseudo random float32 between min and max seeded by the same constance value
func (tc *TemplateContext) RandomFloat(min, max float32) float32 {
	return min + (max-min)*stableRandom.Float32()
}

//NewTemplateContext returns a pointer to TemplateContext
func NewTemplateContext(req *http.Request) *TemplateContext {
	return &TemplateContext{
		Method:     &req.Method,
		Protocol:   &req.Proto,
		Host:       &req.Host,
		RemoteAddr: &req.RemoteAddr,
		Headers:    req.Header,
		FormData:   req.Form,
		PathVars:   pathVarsOrEmpty(req),
		uuid:       make([]byte, 16), // 16 bytes for UUID
	}
}

func pathVarsOrEmpty(req *http.Request) map[string]string {
	pathVars := mux.Vars(req)

	if pathVars == nil {
		// return empty map
		return make(map[string]string)
	}

	return pathVars
}
