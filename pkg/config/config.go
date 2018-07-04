package config

import (
	"errors"
)

type Config struct {
	Backends []Backend
	Policies []Policy
}

type Backend struct {
	Name  string
	Type  string
	Value interface{}
}

type Policy struct {
	Name    string
	Verify  []Probe
	Produce []Producer
}

type Probe struct {
	Name  string
	Type  string
	Value interface{}
}

type Producer struct {
	Name  string
	Type  string
	Value interface{}
}

var (
	ErrBadBackend  = errors.New("bad backend")
	ErrBadProbe    = errors.New("bad probe")
	ErrBadProducer = errors.New("bad producer")
)

func (b *Backend) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v map[string]interface{}
	err := unmarshal(&v)
	if err != nil {
		return err
	}

	if name, ok := v["name"]; ok {
		b.Name = name.(string)
		delete(v, "name")
	}

	if len(v) != 1 {
		return ErrBadBackend
	}

	for type_, value := range v {
		b.Type = type_
		b.Value = value
	}

	return nil
}
func (p *Probe) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v map[string]interface{}
	err := unmarshal(&v)
	if err != nil {
		return err
	}

	if name, ok := v["name"]; ok {
		p.Name = name.(string)
		delete(v, "name")
	}

	if len(v) != 1 {
		return ErrBadProbe
	}

	for type_, value := range v {
		p.Type = type_
		p.Value = value
	}

	return nil
}

func (p *Producer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v map[string]interface{}
	err := unmarshal(&v)
	if err != nil {
		return err
	}

	if name, ok := v["name"]; ok {
		p.Name = name.(string)
		delete(v, "name")
	}

	if len(v) != 1 {
		return ErrBadProducer
	}

	for type_, value := range v {
		p.Type = type_
		p.Value = value
	}

	return nil
}
