package spec

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

// ReadSpecFile reads the server spec from the given file
func ReadSpecFile(path string) (*ServerSpec, error) {
	var err error
	var content []byte
	if content, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}
	return ParseServerSpec(content)
}

// ParseServerSpec loads the server spec from content bytes
func ParseServerSpec(content []byte) (*ServerSpec, error) {
	var serverSpec ServerSpec
	var err error
	if err = yaml.Unmarshal(content, &serverSpec); err != nil {
		return nil, err
	}
	return &serverSpec, nil
}
