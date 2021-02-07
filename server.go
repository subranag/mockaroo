package mockaroo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const (
	// write to the file/create if needed/if exists append
	logFileFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	logFileMode = os.FileMode(0755)
)

type MockServer interface {
	Start() error
}

type muxServer struct {
	conf       *Config
	router     *mux.Router
	reqLogFile *os.File
}

type RequestLog struct {
	RequestUri    *string     `json:"uri"`
	Timestamp     *time.Time  `json:"request_time"`
	Headers       http.Header `json:"headers"`
	Method        *string     `json:"method"`
	ContentLength int64       `json:"content_length"`
	RemoteAddr    *string     `json:"remote_addr"`
	QueryValues   *url.Values `json:"query_params"`
}

// take a http request and convert it into loggable entry
func requestLogFromRequest(r *http.Request) *RequestLog {
	q := r.URL.Query()
	t := time.Now().UTC()
	pt := &t
	return &RequestLog{
		RequestUri:    &r.RequestURI,
		Timestamp:     pt,
		Headers:       r.Header,
		Method:        &r.Method,
		ContentLength: r.ContentLength,
		RemoteAddr:    &r.RemoteAddr,
		QueryValues:   &q,
	}
}

func (s *muxServer) Start() error {

	if s.conf == nil {
		return fmt.Errorf("server config is nil cannot start mockaroo")
	}

	// add all the required routes
	s.addRoutes()

	lfp := s.conf.ServerConfig.RequestLogPath
	if lfp != nil {
		// if logging path is configured setup logging
		lf, err := s.setupLogFile()
		if err != nil {
			return fmt.Errorf("error logging to %v: %w", *lfp, err)
		}
		s.reqLogFile = lf

		// the request log file will be closed when when the
		// server shuts down
		defer lf.Close()
	}

	// add all middlewares
	s.router.Use(s.requestLoggingMiddleware)

	// let the router handle all the requests
	http.Handle("/", s.router)

	// start the server
	http.ListenAndServe(*s.conf.ServerConfig.ListenAddr, nil)

	// all kosher and dandy
	return nil
}

func (s *muxServer) setupLogFile() (*os.File, error) {
	lfp := s.conf.ServerConfig.RequestLogPath
	lf, err := os.OpenFile(*lfp, logFileFlag, logFileMode)
	if err != nil {
		return nil, err
	}

	return lf, nil
}

// log all requests to STDOUT and specified log file if configured
func (s *muxServer) requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: add error handling here
		rl, _ := json.MarshalIndent(requestLogFromRequest(r), "", "    ")

		// write to std out
		fmt.Println(string(rl))

		if s.reqLogFile != nil {
			fmt.Fprintf(s.reqLogFile, "%s\n", rl)
		}

		// call next handler
		next.ServeHTTP(w, r)
	})
}

func (s *muxServer) addRoutes() {
	mocks := s.conf.ServerConfig.Mocks
	router := s.router
	for _, mock := range mocks {
		router.HandleFunc(*mock.Request.Path, genHandleFunc(&mock))
	}
}

// generate the handle function for each mock
func genHandleFunc(mock *Mock) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		for key, val := range mock.Response.Headers {
			resp.Header().Add(key, val)
		}

		fmt.Fprintf(resp, "%s", *mock.Response.ResponseBody)
	}
}

func NewServer(conf *Config) MockServer {
	return &muxServer{conf: conf, router: mux.NewRouter()}
}
