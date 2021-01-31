package mockaroo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
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

	return &config, nil
}
