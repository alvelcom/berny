package interpolation

import (
	"errors"

	//aparser "github.com/mattn/anko/parser"
	//avm "github.com/mattn/anko/vm"

	"github.com/alvelcom/redoubt/api"
)

type Env struct {
	Machine *api.MachineInfo
	User    *api.UserInfo
}

type String interface {
	String(e *Env) (string, error)
}

type StringList interface {
	StringList(e *Env) ([]string, error)
}

func StringVar(dst *String, i interface{}) error {
	if i == nil {
		*dst = nil
		return nil
	}

	switch v := i.(type) {
	case string:
		s := constString(v)
		*dst = &s
		return nil
	case int:
		s := constString(string(v))
		*dst = &s
		return nil
	default:
		return errors.New("NewString: bad type")
	}
}

func StringListVar(dst *StringList, i interface{}) error {
	if i == nil {
		*dst = nil
		return nil
	}

	switch v := i.(type) {
	case []interface{}:
		l := make([]String, len(v))
		*dst = stringList(l)
		for j := range v {
			if err := StringVar(&l[j], v[j]); err != nil {
				return err
			}
		}
		return nil
	case string:
		panic("string list is not impemented")
	}
	return nil
}

type constString string
type stringList []String

func (c *constString) String(e *Env) (string, error) {
	return string(*c), nil
}

func (v stringList) StringList(e *Env) ([]string, error) {
	out := make([]string, len(v))
	for i := range v {
		var err error
		if out[i], err = v[i].String(e); err != nil {
			return out, err
		}
	}
	return out, nil
}
