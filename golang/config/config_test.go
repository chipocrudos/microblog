package config

import "testing"

func TestConfig(t *testing.T) {

	if len(Config.MONGO_URI) < 10 {

		t.Error("Can't load a config variables")
	}
}
