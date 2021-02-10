package mockaroo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

const (
	listenAddrField = "listen_addr"
	maxPortNum      = 65353
)

var validVerbs = map[string]interface{}{
	http.MethodGet:     nil,
	http.MethodHead:    nil,
	http.MethodPost:    nil,
	http.MethodPut:     nil,
	http.MethodPatch:   nil,
	http.MethodDelete:  nil,
	http.MethodConnect: nil,
	http.MethodOptions: nil,
	http.MethodTrace:   nil,
}

type Config struct {
	// self reference to the config file path
	// used only in this package
	configFilePath *string

	ServerConfig ServerConf `hcl:"server,block"`
}

func (c *Config) String() string {
	// ideally this sould not error out so ignoring the err is fine
	// if this panics its going to be tricky
	b, _ := json.Marshal(c)
	return string(b)
}

type ServerConf struct {
	ListenAddr       *string `hcl:"listen_addr"`
	SnakeOilCertPath *string `hcl:"snake_oil_cert"`
	RequestLogPath   *string `hcl:"request_log_path"`
	Mocks            []Mock  `hcl:"mock,block"`
}

type Mock struct {
	Name     string    `hcl:"name,label"`
	Request  *Request  `hcl:"request,block"`
	Response *Response `hcl:"response,block"`
}

type Request struct {
	Path           *string `hcl:"path"`
	NormalizedPath string
	PathPrefix     bool              // should this path be a prefix formulated from the Path
	Verb           *string           `hcl:"verb"`
	Headers        map[string]string `hcl:"headers,optional"` // request match headers
	Queries        map[string]string `hcl:"queries,optional"` // query match headers
}

type Response struct {
	ResponseBody *string           `hcl:"response_body"`
	Headers      map[string]string `hcl:"headers,optional"`
}

type InvalidConfigFile struct {
	path    string
	message string
}

func (e *InvalidConfigFile) Error() string {
	return fmt.Sprintf("invalid config file:%s reason:%s", e.path, e.message)
}

func LoadConfig(filePath *string) (*Config, error) {
	if filePath == nil || strings.TrimSpace(*filePath) == "" {
		return nil, &InvalidConfigFile{path: *filePath, message: "nil or empty config file path"}
	}

	var config Config

	if err := hclsimple.DecodeFile(*filePath, nil, &config); err != nil {
		return nil, &InvalidConfigFile{path: *filePath, message: err.Error()}
	}

	// config file looks good reference it
	config.configFilePath = filePath
	if err := config.validateConfig(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) validateConfig() error {
	listenAddrRegex := regexp.MustCompile(`(?P<host>.+):(?P<port>\d+)`)

	sc := c.ServerConfig
	fp := *c.configFilePath

	if sc.ListenAddr == nil || *sc.ListenAddr == "" {
		errMsg := fmt.Sprintf("%s field in file null or empty", listenAddrField)
		return invalidConfErr(fp, errMsg)
	}

	res := listenAddrRegex.FindStringSubmatch(*sc.ListenAddr)

	if len(res) != 3 {
		errMsg := fmt.Sprintf("expected field %s to be \"<server>:<port>\" found \"%s\"", listenAddrField, *sc.ListenAddr)
		return invalidConfErr(fp, errMsg)
	}

	// not worried about err here see regex we match \d+
	port, _ := strconv.Atoi(res[2])
	if port < 0 || port > maxPortNum {
		errMsg := fmt.Sprintf("port numbers can only be 0 < port < %v found %v in %s=%s", maxPortNum, port, listenAddrField, *sc.ListenAddr)
		return invalidConfErr(fp, errMsg)
	}

	mocks := c.ServerConfig.Mocks

	if len(mocks) == 0 {
		return invalidConfErr(fp, "0 mocks configured, configure mocks using mock:{...} block")
	}

	// now validate all mocks
	for i, mock := range mocks {
		if strings.TrimSpace(mock.Name) == "" {
			errMsg := fmt.Sprintf("invalid empty name for block in index %v, please prvide a valid name", i)
			return invalidConfErr(fp, errMsg)
		}

		if mock.Request == nil {
			errMsg := fmt.Sprintf("request section missing for mock \"%s\"", mock.Name)
			return invalidConfErr(fp, errMsg)
		}

		if err := validPath(fp, &mock); err != nil {
			return err
		}

		// validate verb
		if mock.Request.Verb == nil || strings.TrimSpace(*mock.Request.Verb) == "" {
			errMsg := fmt.Sprintf("null/missing/empty verb for mock \"%s\" verb can only be (GET|HEAD|POST|PUT|DELETE|CONNECT|OPTIONS|TRACE|PATCH)", mock.Name)
			return invalidConfErr(fp, errMsg)
		}

		if _, present := validVerbs[*mock.Request.Verb]; !present {
			errMsg := fmt.Sprintf("invalid verb \"%v\" for mock \"%s\" verb can only be (GET|HEAD|POST|PUT|DELETE|CONNECT|OPTIONS|TRACE|PATCH)", *mock.Request.Verb, mock.Name)
			return invalidConfErr(fp, errMsg)
		}

		// process headers
		reqHeaders := mock.Request.Headers
		for h, v := range reqHeaders {
			_, err := regexp.Compile(v)
			if err != nil {
				errMsg := fmt.Sprintf("invalid request header regexp %s header:\"%s\" in mock \"%s\"", v, h, mock.Name)
				return invalidConfErr(fp, errMsg)
			}
		}

		// process queries
		reqQueries := mock.Request.Queries
		for h, v := range reqQueries {
			_, err := regexp.Compile(v)
			if err != nil {
				errMsg := fmt.Sprintf("invalid request query regexp %s key:\"%s\" in mock \"%s\"", v, h, mock.Name)
				return invalidConfErr(fp, errMsg)
			}
		}

		if mock.Response == nil {
			errMsg := fmt.Sprintf("request section missing for mock \"%s\"", mock.Name)
			return invalidConfErr(fp, errMsg)
		}
	}

	// all validation passed we are kosher
	return nil
}

func validPath(filePath string, mock *Mock) error {
	path := mock.Request.Path
	if path == nil || strings.TrimSpace(*path) == "" {
		errMsg := fmt.Sprintf("request path cannot be nil/\"\" for mock \"%s\"", mock.Name)
		return invalidConfErr(filePath, errMsg)
	}

	//split the path
	parts := strings.Split(*path, "/")

	// the path does not start with a slash it is an error
	if parts[0] != "" {
		errMsg := fmt.Sprintf("request path starts with:\"%v\" anot not \"/\" for mock \"%s\"", parts[0], mock.Name)
		return invalidConfErr(filePath, errMsg)
	}

	for i := 1; i < len(parts); i++ {
		part := parts[i]
		switch {
		case strings.TrimSpace(part) == "":
			if i+1 != len(parts) {
				errMsg := fmt.Sprintf("empty path element path \"%v\" \n", *path)
				errMsg = fmt.Sprintf("%s \" \" white space or empty string cannot be in path; mock is \"%s\"", errMsg, mock.Name)
				return invalidConfErr(filePath, errMsg)
			}
		case strings.Contains(part, "**"):
			if part != "**" || i+1 != len(parts) {
				errMsg := fmt.Sprintf("bad path element \"%v\" in path \"%v\" \n", part, *path)
				errMsg = fmt.Sprintf("%s \"**\" should occur as it is and only at the end of the path for mock \"%s\"", errMsg, mock.Name)
				return invalidConfErr(filePath, errMsg)
			}
			parts[i] = ""
			mock.Request.PathPrefix = true // this path contains path prefix
		case strings.Contains(part, "*"):
			if part != "*" {
				errMsg := fmt.Sprintf("bad path element \"%v\" in path \"%v\" \n", part, *path)
				errMsg = fmt.Sprintf("%s \"*\" should occur as it is; mock is \"%s\"", errMsg, mock.Name)
				return invalidConfErr(filePath, errMsg)
			}
			// all looks good make sure we substitute a variable
			parts[i] = fmt.Sprintf("{pvar%v}", i)
		case strings.Contains(part, "{") || strings.Contains(part, "}"):
			varMatchRegexp := regexp.MustCompile(`^\{.+\}$`)
			if !varMatchRegexp.MatchString(part) {
				errMsg := fmt.Sprintf("bad path element \"%v\" in path \"%v\" \n", part, *path)
				errMsg = fmt.Sprintf("%s variable names should be of form \"{name}\"; mock is \"%s\"", errMsg, mock.Name)
				return invalidConfErr(filePath, errMsg)
			}
		default:
			// all looks good
		}
	}

	// extract and set the normalized path
	mock.Request.NormalizedPath = strings.Join(parts, "/")

	//path looks good
	return nil
}

func invalidConfErr(filPath, message string) error {
	return &InvalidConfigFile{path: filPath, message: message}
}
