package testhelpers

import (
	"os"
	"testing"
)

func TestUnsetEnv(t *testing.T) {
	_ = os.Setenv("ES_URL", "http://previous.env:8080")
	unsetFunc := UnsetEnv("ES_")
	_ = os.Setenv("ES_URL", "http://inner.env:9191")
	unsetFunc()
	if os.Getenv("ES_URL") != "http://previous.env:8080" {
		t.Fail()
	}
}
