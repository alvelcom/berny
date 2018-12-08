package producers

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"net"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"

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
	CommonName hcl.Expression `hcl:"common_name"`
	AltDNS     hcl.Expression `hcl:"alt_dns,optional"`
	AltIPs     hcl.Expression `hcl:"alt_ips,optional"`
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

	commonName, diags := p.CommonName.Value(c.EvalContext)
	if len(diags) > 0 {
		return nil, diags
	}

	if !commonName.Type().Equals(cty.String) {
		return nil, errors.New("producer: common name is not a string")
	}

	altDNS, err := evalStringList(p.AltDNS, c.EvalContext)
	if err != nil {
		return nil, errors.New("producer: alt_dns: " + err.Error())
	}

	altIPsStrings, err := evalStringList(p.AltIPs, c.EvalContext)
	if err != nil {
		return nil, errors.New("producer: alt_ips: " + err.Error())
	}
	var altIPs []net.IP
	for _, ipString := range altIPsStrings {
		altIPs = append(altIPs, net.ParseIP(ipString))
	}

	publicKey := ecdsaKeyResp.PublicKey()
	cert, chain, err := b.Sign(&x509.Certificate{
		Subject: pkix.Name{
			CommonName: commonName.AsString(),
		},
		DNSNames:     altDNS,
		IPAddresses:  altIPs,
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

func evalStringList(expr hcl.Expression, ctx *hcl.EvalContext) ([]string, error) {
	if expr == nil {
		return nil, nil
	}

	evaluated, diags := expr.Value(ctx)
	if len(diags) > 0 {
		return nil, diags
	}

	if evaluated.IsNull() {
		return nil, nil
	}

	if evaluated.Type().IsTupleType() {
		var err error
		evaluated, err = convertTupleToList(cty.String, evaluated)
		if err != nil {
			return nil, err
		}
	}

	if !evaluated.Type().Equals(cty.List(cty.String)) {
		return nil, errors.New("expected a list of strings, got " + evaluated.Type().GoString())
	}

	var list []string
	for _, value := range evaluated.AsValueSlice() {
		list = append(list, value.AsString())
	}

	return list, nil
}

func convertTupleToList(type_ cty.Type, tuple cty.Value) (cty.Value, error) {
	if tuple.LengthInt() == 0 {
		return cty.ListValEmpty(type_), nil
	}

	list := tuple.AsValueSlice()
	for i := range list {
		if !type_.Equals(list[i].Type()) {
			return cty.ListValEmpty(type_), errors.New("can't can tuple to list")
		}
	}
	return cty.ListVal(list), nil
}
