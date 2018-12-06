package producers

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"

	"github.com/alvelcom/redoubt/pkg/api"
	"github.com/alvelcom/redoubt/pkg/backend"
	"github.com/alvelcom/redoubt/pkg/config"
)

var ErrBadProducerType = errors.New("producers: bad type")

type Producer interface {
	Produce(*backend.Map, *hcl.EvalContext) ([]api.Task, []api.Product, error)
}

// PKI
type PKI struct {
	Name       string
	Backend    hcl.Expression `hcl:"backend"`
	CommonName string         `hcl:"common_name"`
	AltDNS     []string       `hcl:"alt_dns,optional"`
	AltIPs     []string       `hcl:"alt_ips,optional"`
}

type File struct {
	Name    string
	Content string `hcl:"content,optional"`
	From    string `hcl:"from,optional"`
}

func New(c config.Producer) (Producer, error) {
	var p Producer
	switch c.Type {
	case "x509":
		p = &PKI{Name: c.Name}
	case "file":
		p = &File{Name: c.Name}
	default:
		return nil, ErrBadProducerType
	}

	diags := gohcl.DecodeBody(c.Config, nil, p)
	if len(diags) > 0 {
		return nil, diags
	}
	return p, nil
}

func (p *PKI) Produce(bm *backend.Map, ec *hcl.EvalContext) ([]api.Task, []api.Product, error) {
	val, diags := p.Backend.Value(ec)
	if len(diags) > 0 {
		return nil, nil, diags
	}

	bn := val.AsString()
	b, ok := bm.X509[bn]
	if !ok {
		return nil, nil, errors.New("no such backend")
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	cert, chain, err := b.Sign(&x509.Certificate{
		Subject: pkix.Name{
			CommonName: "Hey MVP",
		},
		DNSNames:           []string{"example.com"},
		SerialNumber:       big.NewInt(1),
		PublicKey:          &key.PublicKey,
		PublicKeyAlgorithm: x509.ECDSA,
	})
	if err != nil {
		return nil, nil, err
	}

	pemCert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	var pemChain bytes.Buffer
	for i := range chain {
		pem.Encode(&pemChain, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: chain[i],
		})
	}

	ps := []api.Product{
		{
			Name: []string{p.Name, "cert.pem"},
			Body: pemCert,
			Mask: 0400,
		},
		{
			Name: []string{p.Name, "chain.pem"},
			Body: pemChain.Bytes(),
			Mask: 0400,
		},
	}
	return nil, ps, nil
}

func (f *File) Produce(bm *backend.Map, ec *hcl.EvalContext) ([]api.Task, []api.Product, error) {
	content, err := ioutil.ReadFile(f.From)
	if err != nil {
		return nil, nil, err
	}

	ps := []api.Product{{
		Name: []string{f.Name},
		Body: content,
		Mask: 0400,
	}}

	return nil, ps, nil
}
