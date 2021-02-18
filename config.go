package mockaroo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl/v2/hclsimple"
	log "github.com/sirupsen/logrus"
)

type ServerMode int

const (
	HTTP ServerMode = iota
	HTTPS
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

//Config is the root config object that holds entire mockaroo config
type Config struct {
	// self reference to the config file path
	// used only in this package
	configFilePath *string

	ServerConfig *ServerConf `hcl:"server,block"`
}

func (c *Config) String() string {
	// ideally this sould not error out so ignoring the err is fine
	// if this panics its going to be tricky
	b, _ := json.Marshal(c)
	return string(b)
}

//ServerConf mockaroo server configuration
type ServerConf struct {
	ListenAddr       *string `hcl:"listen_addr"`
	SnakeOilCertPath *string `hcl:"snake_oil_cert"`
	SnakeOilKeyPath  *string `hcl:"snake_oil_key"`
	RequestLogPath   *string `hcl:"request_log_path"`
	Mocks            []*Mock `hcl:"mock,block"`
	Mode             ServerMode
}

//Mock matches a specific request and lays out how to generate a response
//to the request
type Mock struct {
	Name     string    `hcl:"name,label"`
	Request  *Request  `hcl:"request,block"`
	Response *Response `hcl:"response,block"`
}

//Request encapsulates a mock request with all information to match a specific
//request
type Request struct {
	Path           *string `hcl:"path"`
	NormalizedPath string
	PathPrefix     bool              // should this path be a prefix formulated from the Path
	Verb           *string           `hcl:"verb"`
	Headers        map[string]string `hcl:"headers,optional"` // request match headers
	Queries        map[string]string `hcl:"queries,optional"` // query match headers
}

//Response encapsulates a complete mock response to a mock Request
type Response struct {
	Staus        int               `hcl:"status,optional"`
	ResponseBody *string           `hcl:"body"`
	ResponseFile *string           `hcl:"file"`
	Headers      map[string]string `hcl:"headers,optional"`
	Delay        *Delay            `hcl:"delay,block"`
	Template     *template.Template
	Content      []byte
}

type Delay struct {
	MaxMillis int64 `hcl:"max_millis"`
	MinMillis int64 `hcl:"min_millis"`
}

//InvalidConfigFile error is raised when given input hcl file fails validation
type InvalidConfigFile struct {
	path    string
	message string
}

func (e *InvalidConfigFile) Error() string {
	return fmt.Sprintf("invalid config file:%s reason:%s", e.path, e.message)
}

//LoadConfig loads the given config file in path and returns a pointed to Config object
//if successful other wise returns a InvalidConfigFile error
func LoadConfig(filePath *string) (*Config, error) {

	if filePath == nil {
		return nil, &InvalidConfigFile{path: "", message: "nil config file path"}
	}

	if strings.TrimSpace(*filePath) == "" {
		return nil, &InvalidConfigFile{path: *filePath, message: "empty config file path"}
	}

	log.Infof("config file : \"%v\"", *filePath)

	var config Config

	if err := hclsimple.DecodeFile(*filePath, nil, &config); err != nil {
		return nil, &InvalidConfigFile{path: *filePath, message: err.Error()}
	}

	log.Info("config file parsed about to validate...")

	// config file parsed
	config.configFilePath = filePath

	// all logical validation
	if err := config.validateConfig(); err != nil {
		return nil, err
	}

	return &config, nil
}

// ALL UN-EXPORTED METHODS

//validateConfig validate the root config object
func (c *Config) validateConfig() error {

	fp := *c.configFilePath

	if c.ServerConfig == nil {
		return invalidConfErr(fp, "server config missing from file")
	}
	sc := c.ServerConfig

	listenAddrRegex := regexp.MustCompile(`(?P<host>.+):(?P<port>\d+)`)
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
	log.Infof("will start server in address: %v", *sc.ListenAddr)

	if c.ServerConfig.RequestLogPath == nil || strings.TrimSpace(*c.ServerConfig.RequestLogPath) == "" {
		// create a temp file and defer closing it but get the path
		tmpfile, err := ioutil.TempFile("", "mockaroo")
		if err != nil {
			return invalidConfErr(fp, fmt.Sprintf("request_log_path not provided and cannot create temp file reason:%v", err.Error()))
		}
		defer tmpfile.Close()
		tmpfileName := tmpfile.Name()
		c.ServerConfig.RequestLogPath = &tmpfileName
	}

	// if key && cert are present then we can start in HTTPS mode
	bothPresent := sc.SnakeOilCertPath != nil && sc.SnakeOilKeyPath != nil

	sc.Mode = HTTP
	log.Info("assuming default mode HTTP")
	if bothPresent {
		sc.Mode = HTTPS
		log.Info("snake oil cert && key present will start in HTTPS mode")
	}

	mocks := c.ServerConfig.Mocks

	if len(mocks) == 0 {
		return invalidConfErr(fp, "0 mocks configured, configure mocks using mock:{...} block")
	}

	// name map to suss out duplicates
	nameToIndex := make(map[string]int)

	// now validate all mocks
	for i, mock := range mocks {
		name := strings.TrimSpace(mock.Name)
		if name == "" {
			errMsg := fmt.Sprintf("invalid empty name for block in index %v, please prvide a valid name", i)
			return invalidConfErr(fp, errMsg)
		}

		prevIndex, present := nameToIndex[name]
		if present {
			errMsg := fmt.Sprintf("mock with name %v already exists in index %v duplicate in %v", name, prevIndex, i)
			return invalidConfErr(fp, errMsg)
		}
		nameToIndex[name] = i

		if mock.Request == nil {
			errMsg := fmt.Sprintf("request section missing for mock \"%s\"", mock.Name)
			return invalidConfErr(fp, errMsg)
		}

		if err := validatePath(fp, mock); err != nil {
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
			errMsg := fmt.Sprintf("response section missing for mock \"%s\"", mock.Name)
			return invalidConfErr(fp, errMsg)
		}

		// if the response Status is set not present or set to 0
		// just assume the response code is going to be success
		if mock.Response.Staus == 0 {
			mock.Response.Staus = 200
		}

		inValidRange := mock.Response.Staus >= 100 && mock.Response.Staus <= 599
		// not in valid range
		if !inValidRange {
			errMsg := fmt.Sprintf("status code is %v, shoud be 100 <= status <= 599 for mock \"%s\"", mock.Response.Staus, mock.Name)
			return invalidConfErr(fp, errMsg)
		}

		if mock.Response.ResponseBody == nil && mock.Response.ResponseFile == nil {
			errMsg := fmt.Sprintf("response section missing body/file atleast one should be present for \"%s\"", mock.Name)
			return invalidConfErr(fp, errMsg)
		}

		if mock.Response.ResponseBody != nil {
			tmplt, err := template.New(mock.Name).Parse(*mock.Response.ResponseBody)
			if err != nil {
				errMsg := fmt.Sprintf("error parsing template for mock \"%s\" error:%s", mock.Name, err.Error())
				return invalidConfErr(fp, errMsg)
			}
			mock.Response.Template = tmplt
		}

		if mock.Response.ResponseFile != nil {
			content, err := ioutil.ReadFile(*mock.Response.ResponseFile)
			if err != nil {
				errMsg := fmt.Sprintf("error reading content from:%v for mock \"%s\" error:%s", *mock.Response.ResponseFile, mock.Name, err.Error())
				return invalidConfErr(fp, errMsg)
			}
			mock.Response.Content = content
		}

		// validate delay
		if mock.Response.Delay != nil {
			minDelay := mock.Response.Delay.MinMillis
			maxDelay := mock.Response.Delay.MaxMillis

			if minDelay < 0 || maxDelay < 0 || maxDelay < minDelay {
				errMsg := fmt.Sprintf("delay min_millis, max_millis >= 0 min_millis <= max_millis and for mock \"%s\" ", mock.Name)
				errMsg = fmt.Sprintf("%s found min_millis:%v max_millis:%v", errMsg, minDelay, maxDelay)
				return invalidConfErr(fp, errMsg)
			}
		}

		// mock looks good
		log.Infof("mock:\"%v\" with path:\"%v\" validates successfully", mock.Name, *mock.Request.Path)
	}

	// all validation passed we are kosher
	return nil
}

//validatePath validate the path of every mock
func validatePath(filePath string, mock *Mock) error {
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
