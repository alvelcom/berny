package backend

import (
	"crypto/x509"
	"errors"

	"github.com/hashicorp/hcl2/gohcl"

	"github.com/alvelcom/redoubt/pkg/config"
)

type Backend interface {
	Type() string
	Init() error
}

type X509 interface {
	Sign(role string, csr *x509.CertificateRequest) ([]x509.Certificate, error)
}

func New(c config.Backend) (Backend, error) {
	var ret Backend
	switch c.Type {
	case "x509_file":
		ret = new(x509File)
	default:
		return nil, errors.New("backend: no backend found")
	}

	diags := gohcl.DecodeBody(c.Config, nil, ret)
	if len(diags) > 0 {
		return nil, diags
	}
	return ret, nil
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
