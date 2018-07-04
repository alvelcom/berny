package backend

import (
	"crypto/x509"
	"errors"

	"github.com/alvelcom/redoubt/pkg/config"
	"github.com/alvelcom/redoubt/pkg/inter"
)

type Backend interface {
	Type() string
	Init(env *inter.Env) error
}

type X509 interface {
	Sign(role string, csr *x509.CertificateRequest) ([]x509.Certificate, error)
}

func New(c config.Backend) (Backend, error) {
	m := map[string]func(config.Backend) (Backend, error){
		"static-x509": newStaticX509,
	}
	if f, ok := m[c.Type]; ok {
		return f(c)
	}
	return nil, errors.New("backend: no backend found")
}

type staticX509 struct {
	Key   inter.String
	Cert  inter.String
	Chain inter.String
}

func newStaticX509(c config.Backend) (Backend, error) {
	m, ok := c.Value.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("static X.509: can't cast to map")
	}

	s := new(staticX509)
	if err := inter.StringVar(&s.Key, m["key"]); err != nil {
		return nil, err
	}
	if err := inter.StringVar(&s.Cert, m["cert"]); err != nil {
		return nil, err
	}
	if err := inter.StringVar(&s.Chain, m["chain"]); err != nil {
		return nil, err
	}

	if s.Key == nil {
		return nil, errors.New("static X.509: key is not set")
	}
	if s.Cert == nil {
		return nil, errors.New("static X.509: cert is not set")
	}

	return s, nil
}

func (s *staticX509) Type() string {
	return "Static X.509"
}
func (s *staticX509) Init(env *inter.Env) error {
	return nil
}

func (s *staticX509) Sign(role string, csr *x509.CertificateRequest) ([]x509.Certificate, error) {
	return nil, nil
}
