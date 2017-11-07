package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/alvelcom/redout/config"
	"github.com/alvelcom/redout/probes"
	"github.com/alvelcom/redout/producers"
)

func init() {
	probes.Register(config.ProbeCasts)
	producers.Register(config.ProducerCasts)
}

func main() {
	data, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		log.Fatal("Can't read config file:", err)
	}

	var c config.Config
	err = yaml.UnmarshalStrict(data, &c)
	if err != nil {
		log.Fatal("Can't unmarshal config file: ", err)
	}

	log.Printf("Config: %+v", c)
}
