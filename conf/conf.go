package conf

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-yaml/yaml"
)

type serverConf struct {
	Addr string `yaml:"addr"`
	Port string
	DSN  string `yaml:"dsn"`
}

var Server serverConf

func init() {
	data, err := ioutil.ReadFile("conf/conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(data, &Server); err != nil {
		log.Fatal(err)
	}
	Server.Port = strings.SplitN(Server.Addr, ":", 2)[1]
}
