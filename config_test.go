package mockaroo

import "testing"

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
