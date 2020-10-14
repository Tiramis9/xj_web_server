package config

import "testing"

func TestInitConfig(t *testing.T) {
	InitConfig("/../config/config.yml")
	t.Log(GetService())
}
