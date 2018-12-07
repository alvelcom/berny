package producers

import (
	"bytes"
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
	"github.com/alvelcom/redoubt/pkg/task"
)

var ErrBadProducerType = errors.New("producers: bad type")

type Context struct {
	Backends      *backend.Map
	TaskResponses TaskResponses
	EvalContext   *hcl.EvalContext
}

type TaskRequests map[[4]string]task.Task
type TaskResponses map[[4]string]task.Response

type Producer interface {
	Prepare(c *Context) (TaskRequests, error)
	Produce(c *Context) ([]api.Product, error)
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

// PKI
type PKI struct {
	Name string

	Backend    hcl.Expression `hcl:"backend"`
	CommonName string         `hcl:"common_name"`
	AltDNS     []string       `hcl:"alt_dns,optional"`
	AltIPs     []string       `hcl:"alt_ips,optional"`
}

func (p *PKI) Prepare(c *Context) (TaskRequests, error) {
	_, ok := c.TaskResponses[[4]string{p.Name}]
	if ok {
		return nil, nil
	}

	return TaskRequests{
		[4]string{p.Name}: &task.ECDSAKey{
			Curve: "P-521",
			Template: api.Product{
				Name: []string{p.Name, "key.pem"},
				Mask: 0600,
			},
		},
	}, nil
}

func (p *PKI) Produce(c *Context) ([]api.Product, error) {
	resp, ok := c.TaskResponses[[4]string{p.Name}]
	if !ok {
		return nil, errors.New("producer: no task response")
	}

	ecdsaKeyResp, ok := resp.(*task.ECDSAKeyResponse)
	if !ok {
		return nil, errors.New("producer: can't cast a task response")
	}

	val, diags := p.Backend.Value(c.EvalContext)
	if len(diags) > 0 {
		return nil, diags
	}

	bn := val.AsString()
	b, ok := c.Backends.X509[bn]
	if !ok {
		return nil, errors.New("no such backend")
	}

	publicKey := ecdsaKeyResp.PublicKey()
	cert, chain, err := b.Sign(&x509.Certificate{
		Subject: pkix.Name{
			CommonName: "Hey MVP",
		},
		DNSNames:     []string{"example.com"},
		SerialNumber: big.NewInt(1),
		PublicKey:    &publicKey,
	})
	if err != nil {
		return nil, err
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
			Mask: 0644,
		},
		{
			Name: []string{p.Name, "chain.pem"},
			Body: pemChain.Bytes(),
			Mask: 0644,
		},
		{
			Name: []string{p.Name, "fullchain.pem"},
			Body: append(pemCert, pemChain.Bytes()...),
			Mask: 0644,
		},
	}
	return ps, nil
}

type File struct {
	Name string

	Content string `hcl:"content,optional"`
	From    string `hcl:"from,optional"`
}

func (f *File) Prepare(c *Context) (TaskRequests, error) {
	return nil, nil
}

func (f *File) Produce(c *Context) ([]api.Product, error) {
	content, err := ioutil.ReadFile(f.From)
	if err != nil {
		return nil, err
	}

	ps := []api.Product{{
		Name: []string{f.Name},
		Body: content,
		Mask: 0600,
	}}

	return ps, nil
}
