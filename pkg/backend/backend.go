package backend

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"

	"github.com/alvelcom/redoubt/pkg/config"
)

type Map struct {
	X509 map[string]X509
}

type X509 interface {
	Sign(cert *x509.Certificate) (derCert []byte, derChain [][]byte, err error)
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

func (x *x509File) Sign(template *x509.Certificate) ([]byte, [][]byte, error) {
	key, err := loadKeyFile(x.Key)
	if err != nil {
		return nil, nil, err
	}

	cert, certDer, err := loadCertFile(x.Cert)
	if err != nil {
		return nil, nil, err
	}

	var chainDer [][]byte
	if len(x.Chain) > 0 {
		chainDer, err = loadChainFile(x.Chain)
		if err != nil {
			return nil, nil, err
		}
	}

	newCert, err := x509.CreateCertificate(rand.Reader, template, cert, template.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	return newCert, append(chainDer, certDer), nil
}

func loadKeyFile(fn string) (*ecdsa.PrivateKey, error) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("backend: can't decode pem")
	}

	return x509.ParseECPrivateKey(block.Bytes)
}

func loadCertFile(fn string) (*x509.Certificate, []byte, error) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, nil, err
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, nil, errors.New("backend: can't decode pem")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, block.Bytes, nil
}

func loadChainFile(fn string) ([][]byte, error) {
	rest, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var certs [][]byte
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if len(rest) == 0 {
			break
		}

		if block == nil {
			return nil, errors.New("backend: can't decode chain's pem")
		}

		if block.Type != "CERTIFICATE" {
			return nil, errors.New("backend: x509: bad pem type in chain: " + block.Type)
		}

		certs = append(certs, block.Bytes)
	}

	return certs, nil
}
