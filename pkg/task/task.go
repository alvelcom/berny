package task

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"

	"github.com/alvelcom/redoubt/pkg/api"
)

type Task interface {
	ToAPI(name []string) api.Task
	Solve() ([]api.Product, Response, error)
}

type Response interface {
	ToAPI(name []string) api.TaskResponse
}

const (
	// Complete list of task types
	ecdsaKeyType = "ecdsa-key"
)

var (
	ErrBadType = errors.New("task: unknown type")
)

func Solve(t api.Task) ([]api.Product, api.TaskResponse, error) {
	task, err := FromAPI(t)
	if err != nil {
		return nil, api.TaskResponse{}, err
	}

	products, resp, err := task.Solve()
	if err != nil {
		return nil, api.TaskResponse{}, err
	}

	return products, resp.ToAPI(t.Name), nil
}

func FromAPI(t api.Task) (Task, error) {
	var task Task
	switch t.Type {
	case ecdsaKeyType:
		task = new(ECDSAKey)
	default:
		return nil, ErrBadType
	}

	err := json.Unmarshal(t.Body, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func FromAPIResponse(r api.TaskResponse) (Response, error) {
	var resp Response
	switch r.Type {
	case ecdsaKeyType:
		resp = new(ECDSAKeyResponse)
	default:
		return nil, ErrBadType
	}

	err := json.Unmarshal(r.Body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ECDSAKey struct {
	Curve    string      `json:"curve"`
	Template api.Product `json:"template"`
}

type ECDSAKeyResponse struct {
	Curve string   `json:"curve"`
	X     *big.Int `json:"x"`
	Y     *big.Int `json:"y"`
}

func (ek ECDSAKey) ToAPI(name []string) api.Task {
	body, err := json.Marshal(ek)
	if err != nil {
		panic(err)
	}

	return api.Task{
		Name: name,
		Type: ecdsaKeyType,
		Body: json.RawMessage(body),
	}
}

func (ek ECDSAKey) Solve() ([]api.Product, Response, error) {
	curve := ECDSACurve(ek.Curve)
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	derKey, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derKey,
	})

	product := ek.Template
	product.Body = pemKey

	products := []api.Product{product}
	resp := ECDSAKeyResponse{
		Curve: ek.Curve,
		X:     key.X,
		Y:     key.Y,
	}
	return products, resp, nil

}

func (er ECDSAKeyResponse) ToAPI(name []string) api.TaskResponse {
	body, err := json.Marshal(er)
	if err != nil {
		panic(err)
	}

	return api.TaskResponse{
		Name: name,
		Type: ecdsaKeyType,
		Body: json.RawMessage(body),
	}
}

func ECDSACurve(c string) elliptic.Curve {
	switch c {
	case "P-224":
		return elliptic.P224()
	case "P-256":
		return elliptic.P256()
	case "P-384":
		return elliptic.P384()
	case "P-521":
		return elliptic.P521()
	default:
		panic("bad curve: " + c)
	}
}
