package mockaroo

import (
	"io/ioutil"
	"os"
	"testing"
)

func configHarness(t *testing.T, configContent string, callback func(string)) {
	f, err := ioutil.TempFile("", "config_test*.hcl")
	if err != nil {
		t.Errorf("failed to open temp file for testing")
	}

	defer f.Close()
	defer os.Remove(f.Name())

	f.WriteString(configContent)
	f.Sync()
	name := f.Name()

	callback(name)
}

func TestNullConfigFails(t *testing.T) {
	conf, err := LoadConfig(nil)

	if err == nil {
		t.Error("error expected for nil config but found to be nil")
	}

	if conf != nil {
		t.Errorf("conf expected to be nil but was found to be non nil:%v", *conf)
	}
}

func TestEmptyConfigFails(t *testing.T) {
	emptyFilePath := "          "
	conf, err := LoadConfig(&emptyFilePath)

	if err == nil {
		t.Error("error expected for nil config but found to be nil")
	}

	if conf != nil {
		t.Errorf("conf expected to be nil but was found to be non nil:%v", *conf)
	}
}

func TestBasicConfigLoadsCorrectly(t *testing.T) {
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

		if *conf.ServerConfig.ListenAddr != "localhost:5000" {
			t.Errorf("did not find the expected ListenAddr:%v found %v", "localhost:5000", *conf.ServerConfig.ListenAddr)
		}

		if len(conf.ServerConfig.Mocks) != 1 {
			t.Errorf("did not find the expected mock count:%v found %v", 1, len(conf.ServerConfig.Mocks))
		}

	})
}

func TestWhenMockNotPresentConfigLoadFails(t *testing.T) {
	sampleConfig := `
	server {
		listen_addr = "localhost:5000"
	}
	`
	configHarness(t, sampleConfig, func(configPath string) {
		_, err := LoadConfig(&configPath)

		if err == nil {
			t.Errorf("expecting config load to fail but config load succeeded")
		}

		assertInvalidConfigError(t, err)
	})
}

func TestMockFailsWhenRequestNotPresent(t *testing.T) {
	sampleConfig := `
	server {
		listen_addr = "localhost:5000"
		mock "hello_world" {
		}
	}
	`
	configHarness(t, sampleConfig, func(configPath string) {
		_, err := LoadConfig(&configPath)

		if err == nil {
			t.Errorf("expecting config load to fail but config load succeeded")
		}

		assertInvalidConfigError(t, err)
	})
}

func TestMockFailsWhenResponseNotPresent(t *testing.T) {
	sampleConfig := `
	server {
		listen_addr = "localhost:5000"
		mock "hello_world" {
			request {
				path = "/hello"
				verb = "GET"
			}
		}
	}
	`
	configHarness(t, sampleConfig, func(configPath string) {
		_, err := LoadConfig(&configPath)

		if err == nil {
			t.Errorf("expecting config load to fail but config load succeeded")
		}

		assertInvalidConfigError(t, err)
	})
}

func assertInvalidConfigError(t *testing.T, err error) {
	if err, ok := err.(*InvalidConfigFile); !ok {
		t.Errorf("the error is not of valid type expected:*InvalidConfigFile found %T", err)
	}
}
