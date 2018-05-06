package probes

import (
	"errors"

	"github.com/alvelcom/redoubt/pkg/config"
)

var (
	ErrBadType = errors.New("probes: bad type")
)

type Probe interface {
	Type() string
}

func New(c config.Probe) (Probe, error) {
	switch c.Type {
	case "gcp":
		return newGCP(c)
	default:
		return nil, ErrBadType
	}
}

type gcp struct {
}

func newGCP(c config.Probe) (Probe, error) {
	return &gcp{}, nil
}

func (g *gcp) Type() string {
	return "gcp"
}
