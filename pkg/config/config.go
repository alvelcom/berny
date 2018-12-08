package config

import (
	"errors"

	"github.com/hashicorp/hcl2/hcl"
)

type Config struct {
	Backends []Backend `hcl:"backend,block"`
	Policies []Policy  `hcl:"policy,block"`
}

type Backend struct {
	Kind   string   `hcl:"kind,label"`
	Name   string   `hcl:"name,label"`
	Type   string   `hcl:"type,attr"`
	Config hcl.Body `hcl:",remain"`
}

type Policy struct {
	Name    string     `hcl:"name,label"`
	Verify  []Probe    `hcl:"verify,block"`
	Produce []Producer `hcl:"produce,block"`
}

type Probe struct {
	Type   string   `hcl:"type,label"`
	Config hcl.Body `hcl:",remain"`
}

type Producer struct {
	ProducerHeader
	Type   string   `hcl:"type,label"`
	Name   string   `hcl:"name,label"`
	Config hcl.Body `hcl:",remain"`
}

type ProducerHeader struct {
	Type string `hcl:"type,label"`
	Name string `hcl:"name,label"`
}

var (
	ErrBadBackend  = errors.New("bad backend")
	ErrBadProbe    = errors.New("bad probe")
	ErrBadProducer = errors.New("bad producer")
)
