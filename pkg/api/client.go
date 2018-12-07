package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client interface {
	Harvest(r []TaskResponse) ([]Product, []Task, []Error, error)
}

type HTTPClient struct {
	client       *http.Client
	url          string
	serverCookie string
	info         MachineInfo
}

func NewHTTPClient(c *http.Client, url string, info MachineInfo) (*HTTPClient, error) {
	return &HTTPClient{
		client: c,
		url:    url,
		info:   info,
	}, nil
}

func (hc *HTTPClient) Harvest(r []TaskResponse) (p []Product, t []Task, e []Error, err error) {
	var b bytes.Buffer
	if err = json.NewEncoder(&b).Encode(Request{
		ClientVersion: 0,
		ServerCookie:  hc.serverCookie,
		Machine:       &hc.info,
		TaskResponses: r,
	}); err != nil {
		return
	}

	req, err := http.NewRequest("POST", hc.url+"/v1/harvest", &b)
	if err != nil {
		return
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var answer Response
	if err = json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		return
	}

	p = answer.Products
	t = answer.Tasks
	e = answer.Errors
	err = nil
	return
}
