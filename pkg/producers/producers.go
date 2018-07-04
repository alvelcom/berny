package producers

import (
	"errors"
	"io/ioutil"

	"github.com/alvelcom/redoubt/pkg/api"
	"github.com/alvelcom/redoubt/pkg/config"
	"github.com/alvelcom/redoubt/pkg/inter"
)

var ErrBadProducerType = errors.New("producers: bad type")

type Producer interface {
	Type() string
	Produce(*inter.Env) ([]api.Task, []api.Product, error)
}

// PKI
type PKI struct {
	Backend inter.String
	Profile inter.String
	CRT     CRT
}

type CRT struct {
	CommonName inter.String
	AltDNS     inter.StringList
	AltIPs     inter.StringList
}

type File struct {
	Files inter.StringMap
}

func New(c config.Producer) (Producer, error) {
	switch c.Type {
	case "pki":
		return newPKI(c)
	case "file":
		return newFile(c)
	default:
		return nil, ErrBadProducerType
	}
}

func newPKI(c config.Producer) (Producer, error) {
	params, ok := c.Value.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("new pki: can't cast to map")
	}

	var pki PKI

	if err := inter.StringVar(&pki.Backend, params["backend"]); err != nil {
		return nil, err
	}
	if err := inter.StringVar(&pki.Profile, params["profile"]); err != nil {
		return nil, err
	}

	crt_, ok := params["crt"]
	if !ok {
		return nil, errors.New("new pki: no crt block")
	}
	crt, ok := crt_.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("new pki: can't cast crt to map")
	}

	if err := inter.StringVar(&pki.CRT.CommonName, crt["common-name"]); err != nil {
		return nil, err
	}
	if err := inter.StringListVar(&pki.CRT.AltDNS, crt["alt-dns"]); err != nil {
		return nil, err
	}
	if err := inter.StringListVar(&pki.CRT.AltIPs, crt["alt-ips"]); err != nil {
		return nil, err
	}

	return &pki, nil
}

func (p *PKI) Type() string {
	return "pki"
}

func (p *PKI) Produce(*inter.Env) ([]api.Task, []api.Product, error) {

	return nil, []api.Product{{}}, nil
}

func newFile(c config.Producer) (Producer, error) {
	params, ok := c.Value.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("new pki: can't cast to map")
	}

	var file File
	if err := inter.StringMapVar(&file.Files, params); err != nil {
		return nil, err
	}

	return &file, nil
}

func (p *File) Type() string {
	return "pki"
}

func (p *File) Produce(e *inter.Env) ([]api.Task, []api.Product, error) {
	m, err := p.Files.StringMap(e)
	if err != nil {
		return nil, nil, err
	}

	var ps []api.Product
	for k, v := range m {
		content, err := ioutil.ReadFile(v)
		if err != nil {
			return nil, ps, err
		}

		ps = append(ps, api.Product{
			Name: []string{k},
			Body: content,
		})
	}

	return nil, ps, nil
}
