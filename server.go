package mockaroo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	// gorilla seems like the best fit, supports a lot of rich matching
	// and plays well with native golng http
	"github.com/gorilla/mux"
)

const (
	// write to the file/create if needed/if exists append
	logFileFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	// default mode for log file creation
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
	// if the server fails to start it will return an error
	if s.conf.ServerConfig.Mode == HTTPS {
		return http.ListenAndServeTLS(*s.conf.ServerConfig.ListenAddr,
			*s.conf.ServerConfig.SnakeOilCertPath,
			*s.conf.ServerConfig.SnakeOilKeyPath, nil)
	}
	return http.ListenAndServe(*s.conf.ServerConfig.ListenAddr, nil)
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
		rl, _ := json.Marshal(requestLogFromRequest(r))

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
	for _, v := range s.conf.ServerConfig.Mocks {
		s.router.HandleFunc(*v.Request.Path, genHandleFunc(v))
	}
}

// generate the handle function for each mock
func genHandleFunc(mock Mock) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		for key, val := range mock.Response.Headers {
			resp.Header().Add(key, val)
		}
		if mock.Response.ResponseBody != nil {

		}

		switch {
		case mock.Response.Template != nil:
			// TODO: pass all context data here
			err := mock.Response.Template.Execute(resp, nil)
			if err != nil {
				// raise a 500
				resp.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(resp, "template execution failed for mock \"%v\" error:%v", mock.Name, err.Error())
			}
		case mock.Response.Content != nil:
			fmt.Fprintf(resp, "%s", mock.Response.Content)
		default:
			// we should never be here if we are here mockaroo bunged it
			// please open an issue
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(resp, "BAD BAD BAD mockaroo fix your test cases; mock \"%v\"", mock.Name)
		}
	}
}

// NewServer creates a mock server with the given configuration
func NewServer(conf *Config) MockServer {
	return &muxServer{conf: conf, router: mux.NewRouter()}
}
