package LinkServer

import "testing"

func TestStart(t *testing.T) {
	Start()
}

func TestInitServerFromConfig(t *testing.T) {
	InitServerFromConfig("./config/config.json")
}
