package serve

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/alvelcom/redoubt/config"
	"github.com/alvelcom/redoubt/probes"
	"github.com/alvelcom/redoubt/producers"
)

var fs = flag.NewFlagSet("redoubt serve", flag.ExitOnError)
var (
	listenAddr = fs.String("listen", "0.0.0.0:2326",
		`Listen for incomming request there`)
	configFile = fs.String("config", "test.yaml",
		`Configuration file to use`)
)

func init() {
	probes.Register()
	producers.Register()
}

func Main(args []string) {
	fs.Parse(args)
	log.Print(*listenAddr)

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal("Can't read config file: ", err)
	}

	var c config.Config
	err = yaml.UnmarshalStrict(data, &c)
	if err != nil {
		log.Fatal("Can't unmarshal config file: ", err)
	}

	log.Printf("Config: %+v", c)
}
