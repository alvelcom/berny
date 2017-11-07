package config

import (
	"context"
    "errors"
)

type Config struct {
	Policies []Policy
}

type Policy struct {
	Name    string
	Verify  []Probe
    Produce []Producer
}

type Probe struct {
    Name  string
	Probe interface {
        Type() string
        Test(ctx context.Context) error
    }
}

type Producer struct {
    Name string
    Producer interface {
        Type() string

    }
}

var (
    ErrCantUnmarshal = errors.New("Can't UnmarshalYAML probes or outputs struct")
)

// Casting magic for probes
type ProbeCast func(c *Probe, unmarshal func(interface{}) error) error
var ProbeCasts = map[string]ProbeCast{}

func (c *Probe) UnmarshalYAML(unmarshal func(interface{}) error) error {
    m := map[string]interface{}{}
    unmarshal(m)

    for name, cast := range ProbeCasts {
        if _, ok := m[name]; ok {
            return cast(c, unmarshal)
        }
    }

	return ErrCantUnmarshal
}

// Casting magic for producers
type ProducerCast func(c *Producer, unmarshal func(interface{}) error) error
var ProducerCasts = map[string]ProducerCast{}

func (c *Producer) UnmarshalYAML(unmarshal func(interface{}) error) error {
    m := map[string]interface{}{}
    unmarshal(m)

    for name, cast := range ProducerCasts {
        if _, ok := m[name]; ok {
            return cast(c, unmarshal)
        }
    }

	return ErrCantUnmarshal
}
