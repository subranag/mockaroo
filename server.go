package mockaroo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ContentType = "Content-Type"
	JSON        = "application/json"
)

type MockServer interface {
	Start()
}

type muxServer struct {
	conf   *Config
	router *mux.Router
}

func (s *muxServer) Start() {
	s.router.HandleFunc("/ping", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add(ContentType, JSON)
		pong := map[string]string{"message": "pong"}
		bytes, _ := json.Marshal(pong)
		fmt.Fprintf(resp, "%s", bytes)
	})

	s.addRoutes()
	http.Handle("/", s.router)

	http.ListenAndServe(*s.conf.ServerConfig.ListenAddr, nil)
}

func (s *muxServer) addRoutes() {
	mocks := s.conf.ServerConfig.Mocks
	router := s.router
	for _, mock := range mocks {
		router.HandleFunc(*mock.Path, genHandleFunc(&mock))
	}
}

func genHandleFunc(mock *Mock) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		for key, val := range mock.Headers {
			resp.Header().Add(key, val)
		}

		fmt.Fprintf(resp, "%s", *mock.ResponseBody)
	}
}

func NewServer(conf *Config) MockServer {
	return &muxServer{conf: conf, router: mux.NewRouter()}
}
