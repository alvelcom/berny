package producers

import (
	"errors"
	"io/ioutil"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"

	"github.com/alvelcom/redoubt/pkg/api"
	"github.com/alvelcom/redoubt/pkg/config"
)

var ErrBadProducerType = errors.New("producers: bad type")

type Producer interface {
	Type() string
	Produce(*hcl.EvalContext) ([]api.Task, []api.Product, error)
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

func (p *PKI) Type() string {
	return "pki"
}

func (p *PKI) Produce(*hcl.EvalContext) ([]api.Task, []api.Product, error) {
	return nil, nil, nil
}

func (f *File) Type() string {
	return "file"
}

func (f *File) Produce(e *hcl.EvalContext) ([]api.Task, []api.Product, error) {
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
