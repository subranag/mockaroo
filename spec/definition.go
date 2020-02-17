package spec

import "net/http"

// MockAction defines the action that needs to be taken to process
// a given HTTP request and mock and equvalent reponse
type MockAction interface {
	PerformMockAction(w http.ResponseWriter, r *http.Request)
}

// ServerSpec holds the details about the mock server
type ServerSpec struct {
	// Port on which this mock server needs to run
	Port int `yaml:"server_port"`

	// Mocks all the mock specs for this server
	Mocks []MockSpec `yaml:"mock_specs"`
}

// MockMatcher matches a request and returns true if match false otherwise
// the complete context of the http request is taken into account
type MockMatcher interface {
	MatchMock(r *http.Request) bool
}

// MockSpec defines the mocking spec for mockaroo
// the mock spec is usually writeen in a JSON mock file
// read below for the specifications of the MockFile
type MockSpec struct {
	// HTTPMethod specifies the method to match this mock spec on
	HTTPMethod string `yaml:"http_method"`

	// MatchPath specifies the exact path or path prefix to match on
	// examples include "/student/create" or "/student/create/*"
	// the most specific match is chosen and the action is applied
	MatchPath string `yaml:"path"`

	// MatchRequestParams specifies the list of request parameters to match on
	// you can create different mock actions
	MatchRequestParams []string `yaml:"params_present"`

	// Action is the actual mock action that needs to ber performed for the matched spec
	Action struct {
		// ResponseCode to be sent back in response to this matched MockSpec
		ResponseCode int `yaml:"response_code"`

		// MinDelayMillis the minimun delay in millis before sending the mock response
		// if not present there will be no mock delay
		MinDelayMillis uint64 `yaml:"min_delay_ms"`

		// MaxDelayMillis the maximum delay in millis before sending the mock response
		// if MinDelayMillis are specified and MaxDelayMillis > MinDelayMillis
		// the mock response will be delayed by RANDOM(MinDelayMillis, MaxDelayMillis)
		MaxDelayMillis uint64 `yaml:"max_delay_ms"`

		// ResponseHeaders that need to be sent as a part of this match spec
		ResponseHeaders map[string]string `yaml:"response_headers"`

		// ResponseTemplate that can contain template expressions for returning the response
		ResponseTemplate string `yaml:"response_template"`

		// ResponseFile the file bytes will be read and processed as response
		ResponseFile string `yaml:"response_file"`

		// ResponseString plain response string which is simply echoed back as response
		ResponseString string `yaml:"response_string"`
	} `yaml:"action"`
}

// ResponseModel will be data model passed to the response template
type ResponseModel struct {
	RequestPath string
}
