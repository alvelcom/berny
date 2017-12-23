package producers

import (
	"errors"

	"github.com/alvelcom/redoubt/config"
)

var ErrBadProducerType = errors.New("producers: bad type")

type Producer interface {
	Type() string
}

// PKI
type PKI struct {
	Mount string
	Role  string
	CSR   CSR
}

type CSR struct {
	CommonName string   `yaml:"common-name"`
	AltDNS     []string `yaml:"alt-dns"`
	AltIPs     []string `yaml:"alt-ip"`
}

// Secret
type Secret struct {
}

func New(c config.Producer) (Producer, error) {
	switch c.Type {
	case "pki":
		return newPKI(c)
	default:
		return nil, ErrBadProducerType
	}
}

func newPKI(c config.Producer) (Producer, error) {
	return &PKI{}, nil
}

func (p *PKI) Type() string {
	return "pki"
}
