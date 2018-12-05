package backend

import (
	"crypto/x509"
	"errors"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"

	"github.com/alvelcom/redoubt/pkg/config"
)

type Map struct {
	X509 map[string]X509
}

type X509 interface {
	Sign(role string, csr *x509.CertificateRequest) ([]x509.Certificate, error)
}

func NewMap() *Map {
	return &Map{
		X509: make(map[string]X509),
	}
}

func (m *Map) Add(c config.Backend) error {
	var x509 X509

	switch c.Type {
	case "x509_file":
		x509 = new(x509File)
	default:
		return errors.New("backend: no backend found")
	}

	var diags hcl.Diagnostics
	switch {
	case x509 != nil:
		diags = gohcl.DecodeBody(c.Config, nil, x509)
		m.X509[c.Name] = x509
	default:
		panic("backend: Add: what?")
	}

	if len(diags) > 0 {
		return diags
	}
	return nil
}

type x509File struct {
	Key   string `hcl:"key"`
	Cert  string `hcl:"cert"`
	Chain string `hcl:"chain,optional"`
}

func (s *x509File) Type() string {
	return "Static X.509"
}
func (s *x509File) Init() error {
	return nil
}

func (s *x509File) Sign(role string, csr *x509.CertificateRequest) ([]x509.Certificate, error) {
	return nil, nil
}
