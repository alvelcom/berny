package producers

import (
	"errors"

	"github.com/alvelcom/redoubt/config"
	inter "github.com/alvelcom/redoubt/interpolation"
)

var ErrBadProducerType = errors.New("producers: bad type")

type Producer interface {
	Type() string
}

// PKI
type PKI struct {
	Backend inter.String
	Role    inter.String
	CRT     CRT
}

type CRT struct {
	CommonName inter.String
	AltDNS     inter.StringList
	AltIPs     inter.StringList
}

func New(c config.Producer) (Producer, error) {
	switch c.Type {
	case "pki":
		return newPKI(c)
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
	if err := inter.StringVar(&pki.Role, params["role"]); err != nil {
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
