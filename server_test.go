package linkserver

import "testing"

func TestStart(t *testing.T) {
	Start(8888)
}

func TestInitServerFromConfig(t *testing.T) {
	InitServerFromConfig("./config/config.json")
}
