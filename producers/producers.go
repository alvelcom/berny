package producers

import (
	"context"
	"fmt"

	"github.com/alvelcom/redoubt/config"
)

func Register() {
	config.ProducerCasts["pki"] = castPKI
	config.ProducerCasts["secret"] = castSecret
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

func castPKI(c *config.Producer, unmarshal func(interface{}) error) error {
	var t struct {
		Name string

		PKI *PKI `yaml:"pki"`
	}
	if err := unmarshal(&t); err != nil {
		return err
	}
	if t.PKI == nil {
		return config.ErrCantUnmarshal
	}
	c.Name = t.Name
	c.Producer = t.PKI
	return nil
}

func (p *PKI) Type() string {
	return "pki"
}

func (p *PKI) Test(ctx context.Context) error {
	return nil
}

func (p *PKI) String() string {
	return fmt.Sprintf("%#v", *p)
}

// Secret
type Secret struct {
}

func castSecret(c *config.Producer, unmarshal func(interface{}) error) error {
	var t struct {
		Name string

		Secret *Secret `yaml:"secret"`
	}
	if err := unmarshal(&t); err != nil {
		return err
	}
	if t.Secret == nil {
		return config.ErrCantUnmarshal
	}
	c.Name = t.Name
	c.Producer = t.Secret
	return nil
}

func (p *Secret) Type() string {
	return "secret"
}

func (p *Secret) Test(ctx context.Context) error {
	return nil
}

func (p *Secret) String() string {
	return fmt.Sprintf("%#v", *p)
}
