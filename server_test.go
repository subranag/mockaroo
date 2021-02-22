package mockaroo

import (
	"net/http"
	"net/http/httptest"
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
		conf, err := LoadConfig(&configPath)

		if err != nil {
			t.Errorf("config load failed with error:%v>\n%s", err.Error(), sampleConfig)
		}

		// make new server with config
		muxServer := &muxServer{conf: conf, router: mux.NewRouter()}

		// make test HTTP request
		req, err := http.NewRequest("GET", "/hello", nil)
		if err != nil {
			t.Errorf("cannot create new request error:%v", err)
		}

		// add all routes
		muxServer.addRoutes()

		rr := httptest.NewRecorder()
		muxServer.router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200 but HTTP request failed with:%v", rr.Code)
		}
	})
}
