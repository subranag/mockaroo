package mockaroo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
