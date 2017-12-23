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

type Policy struct {
	Name    string
	Verify  []probes.Probe
	Produce []producers.Producer
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

	policies, err := castPolicies(c.Policies)
	if err != nil {
		log.Fatal("Can't initialize policies: ", err)
	}
	log.Printf("Policies: %#v", policies)
}

func castPolicies(ps []config.Policy) ([]Policy, error) {
	policies := []Policy{}
	for _, p := range ps {
		policy := Policy{
			Name:    p.Name,
			Verify:  []probes.Probe{},
			Produce: []producers.Producer{},
		}
		for _, probe := range p.Verify {
			probe, err := probes.New(probe)
			if err != nil {
				return nil, err
			}
			policy.Verify = append(policy.Verify, probe)
		}
		for _, producer := range p.Produce {
			producer, err := producers.New(producer)
			if err != nil {
				return nil, err
			}
			policy.Produce = append(policy.Produce, producer)
		}
		policies = append(policies, policy)
	}
	return policies, nil
}
