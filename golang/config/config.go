package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Debug     bool   `envconfig:"DEBUG"`
	PORT      string `required:"true" envconfig:"PORT"`
	MONGO_URI string `required:"true" envconfig:"MONGO_URI"`
	MONGO_DB  string `required:"true" envconfig:"MONGO_DB"`
	JWT_SALT  string `required:"true" envconfig:"JWT_SALT"`
	EXP_TIME  int    `required:"true" envconfig:"EXP_TIME"`
}

func (this *Configuration) LoadConfig() {

	err := envconfig.Process("", this)

	if err != nil {
		log.Fatal(err.Error())
	}

}

var Config = new(Configuration)

func init() {
	Config.LoadConfig()
	log.Println("Apply app config")
}
