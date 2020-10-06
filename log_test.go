package bitflyergo

import (
	"testing"
)

func TestLogln(t *testing.T) {
	Logger = nil
	logln("HelloWorld")
}
