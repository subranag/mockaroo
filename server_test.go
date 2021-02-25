package mockaroo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestBasicMockWorksCorrectly(t *testing.T) {
	sampleConfig := `
	server {
		listen_addr = "localhost:5000"
		mock "hello_world" {
			request {
				path = "/hello"
				verb = "GET"
			}
			response {
				body = <<EOF
				world
				EOF
			}
		}
	}
	`
	configHarness(t, sampleConfig, func(configPath string) {

		muxServer := loadConfigAndGetServer(t, configPath, sampleConfig)

		// make test HTTP request
		req, err := http.NewRequest("GET", "/hello", nil)
		if err != nil {
			t.Errorf("cannot create new request error:%v", err)
		}

		rr := httptest.NewRecorder()
		muxServer.router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 but HTTP request failed with:%v", rr.Code)
		}

		responseBody := strings.TrimSpace(string(rr.Body.Bytes()))
		if responseBody != "world" {
			t.Errorf("expected response body:world but found:%v", responseBody)
		}
	})
}

func TestQueryMatchingWorksCorrectly(t *testing.T) {
	sampleConfig := `
	server {
		listen_addr = "localhost:5000"
		mock "query1" {
			request {
				path = "/foo"
				verb = "GET"
				queries = {
					bar = ""
				}
			}
			response {
				body = <<EOF
				foo bar
				EOF
			}
		}

		mock "query2" {
			request {
				path = "/foo"
				verb = "GET"
				queries = {
					bared = ""
				}
			}
			response {
				body = <<EOF
				donkey doo
				EOF
			}
		}
	}
	`
	configHarness(t, sampleConfig, func(configPath string) {

		muxServer := loadConfigAndGetServer(t, configPath, sampleConfig)

		rr := httptest.NewRecorder()
		muxServer.router.ServeHTTP(rr, createGetRequest(t, "/foo?bar", nil))

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 but HTTP request failed with:%v", rr.Code)
		}

		responseBody := strings.TrimSpace(string(rr.Body.Bytes()))
		if responseBody != "foo bar" {
			t.Errorf("expected response body:world but found:%v", responseBody)
		}

		// make sure other match works correctly
		rr = httptest.NewRecorder()
		muxServer.router.ServeHTTP(rr, createGetRequest(t, "/foo?bared", nil))

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 but HTTP request failed with:%v", rr.Code)
		}

		responseBody = strings.TrimSpace(string(rr.Body.Bytes()))
		if responseBody != "donkey doo" {
			t.Errorf("expected response body:donkey doo but found:%v", responseBody)
		}
	})
}

func TestHeaderMatchingWorksCorrectly(t *testing.T) {
	sampleConfig := `
	server {
		listen_addr = "localhost:5000"
		mock "query1" {
			request {
				path = "/foo"
				verb = "GET"
				headers = {
					origin = "safehost"
				}
			}
			response {
				body = <<EOF
				foo bar
				EOF
			}
		}

		mock "query2" {
			request {
				path = "/foo"
				verb = "GET"
				headers = {
					origin = "dodgyhost"
				}
			}
			response {
				body = <<EOF
				donkey doo
				EOF
			}
		}
	}
	`
	configHarness(t, sampleConfig, func(configPath string) {

		muxServer := loadConfigAndGetServer(t, configPath, sampleConfig)

		rr := httptest.NewRecorder()
		headers := map[string]string{"origin": "safehost"}
		muxServer.router.ServeHTTP(rr, createGetRequest(t, "/foo?bar", headers))

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 but HTTP request failed with:%v", rr.Code)
		}

		responseBody := strings.TrimSpace(string(rr.Body.Bytes()))
		if responseBody != "foo bar" {
			t.Errorf("expected response body:world but found:%v", responseBody)
		}

		// make sure other match works correctly
		rr = httptest.NewRecorder()
		headers = map[string]string{"origin": "dodgyhost"}
		muxServer.router.ServeHTTP(rr, createGetRequest(t, "/foo?bared", headers))

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 but HTTP request failed with:%v", rr.Code)
		}

		responseBody = strings.TrimSpace(string(rr.Body.Bytes()))
		if responseBody != "donkey doo" {
			t.Errorf("expected response body:donkey doo but found:%v", responseBody)
		}
	})
}

func TestResponseDelayWorksCorrectly(t *testing.T) {
	sampleConfig := `
	server {
		listen_addr = "localhost:5000"
		mock "hello_world" {
			request {
				path = "/hello"
				verb = "GET"
			}
			response {
				body = <<EOF
				world
				EOF

				headers = {
					foo = "bar"
				}

				delay {
					min_millis = 100
					max_millis = 100
				}
			}
		}
	}
	`
	configHarness(t, sampleConfig, func(configPath string) {

		muxServer := loadConfigAndGetServer(t, configPath, sampleConfig)

		// make test HTTP request
		req, err := http.NewRequest("GET", "/hello", nil)
		if err != nil {
			t.Errorf("cannot create new request error:%v", err)
		}

		rr := httptest.NewRecorder()
		start := time.Now()
		muxServer.router.ServeHTTP(rr, req)
		elapsed := time.Since(start)

		if elapsed.Milliseconds() < 100 {
			t.Errorf("configured delay was 100 millis but delay laster only for:%v", elapsed.Milliseconds())
		}

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 but HTTP request failed with:%v", rr.Code)
		}

		if rr.HeaderMap.Get("foo") != "bar" {
			t.Errorf("expected response header foo to be bar but was:%v", rr.HeaderMap.Get("foo"))
		}

		responseBody := strings.TrimSpace(string(rr.Body.Bytes()))
		if responseBody != "world" {
			t.Errorf("expected response body:world but found:%v", responseBody)
		}
	})
}

func TestRequestLogFromRequestWorksCorrectly(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello?foo=bar", nil)

	if err != nil {
		t.Errorf("cannot create new request error:%v", err)
	}

	reql := requestLogFromRequest(req)

	if reql == nil {
		t.Error("request log expected to be non nil")
	}

	if *reql.Method != "GET" {
		t.Errorf("expected method GET but found %v", *reql.Method)
	}
}

func createGetRequest(t *testing.T, query string, headers map[string]string) *http.Request {
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		t.Errorf("cannot create new request error:%v", err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req
}

func loadConfigAndGetServer(t *testing.T, configPath, sampleConfig string) *muxServer {
	conf, err := LoadConfig(&configPath)

	if err != nil {
		t.Errorf("config load failed with error:%v>\n%s", err.Error(), sampleConfig)
	}

	// make new server with config
	muxServer := &muxServer{conf: conf, router: mux.NewRouter()}

	// add all routes
	muxServer.addRoutes()
	return muxServer
}
