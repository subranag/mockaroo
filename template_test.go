package mockaroo

import (
	"net/http"
	"testing"
)

func TestTemplateContextFunctions(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello?foo=bar", nil)

	if err != nil {
		t.Errorf("cannot create new request error:%v", err)
	}

	tc := NewTemplateContext(req)

	uuid := tc.NewUUID()
	expected := "711580a1-ffff-5c48-9ba4-46c8187f2398"
	// i can do this because of stable pseudo random number generation
	if uuid != expected {
		t.Errorf("expected uuid:%v found:%v", expected, uuid)
	}

	rint := tc.RandomInt(10, 15)
	if rint < 10 || rint > 15 {
		t.Errorf("expected random int to be between 10 and 15 but found:%v", rint)
	}
}
