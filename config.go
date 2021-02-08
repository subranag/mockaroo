package mockaroo

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

const (
	listenAddrField = "listen_addr"
	maxPortNum      = 65353
)

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
	Path *string `hcl:"path"`
	Verb *string `hcl:"verb"`
}

type Response struct {
	ResponseBody *string           `hcl:"response_body"`
	Headers      map[string]string `hcl:"headers"`
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

	if sc.ListenAddr == nil || *sc.ListenAddr == "" {
		return c.invalidConfErr(fmt.Sprintf("%s field in file null or empty", listenAddrField))
	}

	res := listenAddrRegex.FindStringSubmatch(*sc.ListenAddr)

	if len(res) != 3 {
		return c.invalidConfErr(fmt.Sprintf("expected field %s to be \"<server>:<port>\" found \"%s\"", listenAddrField, *sc.ListenAddr))
	}

	// not worried about err here see regex we match \d+
	port, _ := strconv.Atoi(res[2])
	if port < 0 || port > maxPortNum {
		return c.invalidConfErr(fmt.Sprintf("port numbers can only be 0 < port < %v found %v in %s=%s", maxPortNum, port, listenAddrField, *sc.ListenAddr))
	}

	mocks := c.ServerConfig.Mocks

	if len(mocks) == 0 {
		return c.invalidConfErr("0 mocks configured, configure mocks using mock:{...} block")
	}

	// now validate all mocks
	for i, mock := range mocks {
		if strings.TrimSpace(mock.Name) == "" {
			return c.invalidConfErr(fmt.Sprintf("invalid empty name for block in index %v, please prvide a valid name", i))
		}

		if mock.Request == nil {
			return c.invalidConfErr(fmt.Sprintf("request section missing for mock \"%s\"", mock.Name))
		}

		if mock.Request.Path == nil {
			return c.invalidConfErr(fmt.Sprintf("request path cannot be nil for mock \"%s\"", mock.Name))
		}

		pathRegexp := regexp.MustCompile(`/((?:[^/]*/)*)(.*)`)
		res := pathRegexp.FindAllString(*mock.Request.Path, -1)
		if len(res) == 0 {
			return c.invalidConfErr(fmt.Sprintf("not a valid path:%s for mock \"%s\"", *mock.Request.Path, mock.Name))
		}

		improperStar := regexp.MustCompile(`/\w+\*|/\*\w+|\s\*|\*\s`)
		res = improperStar.FindAllString(*mock.Request.Path, -1)
		if len(res) > 0 {
			return c.invalidConfErr(fmt.Sprintf("\"*\" charecter present in improper place in path:%s for mock \"%s\"", *mock.Request.Path, mock.Name))
		}

		prefixPathRegexp := regexp.MustCompile(`/\*+`)
		res = prefixPathRegexp.FindAllString(*mock.Request.Path, -1)
		fmt.Printf("%v\n", res)

		if mock.Response == nil {
			return c.invalidConfErr(fmt.Sprintf("request section missing for mock \"%s\"", mock.Name))
		}
	}

	// all validation passed we are kosher
	return nil
}

func (c *Config) invalidConfErr(message string) error {
	return &InvalidConfigFile{path: *c.configFilePath, message: message}
}
