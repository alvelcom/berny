package probes

import (
	"context"
	"fmt"

	"github.com/alvelcom/redoubt/config"
)

func Register() {
	config.ProbeCasts["gce"] = castGCE
}

// GCE
type GCE struct {
}

func castGCE(c *config.Probe, unmarshal func(interface{}) error) error {
	var t struct {
		Name string
		GCE  *GCE `yaml:"gce"`
	}
	if err := unmarshal(&t); err != nil {
		return err
	}
	c.Name = t.Name
	c.Probe = t.GCE
	return nil
}

func (p *GCE) Type() string {
	return "gce"
}

func (p *GCE) Test(ctx context.Context) error {
	return nil
}

func (p *GCE) String() string {
	return fmt.Sprintf("%#v", *p)
}
