package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client interface {
	Harvest(r []Returning) ([]Product, []Task, []Error, error)
}

type httpClient struct {
	client       *http.Client
	url          string
	serverCookie string
	info         MachineInfo
}

func NewClient(c *http.Client, url string, info MachineInfo) (Client, error) {
	return &httpClient{
		client: c,
		url:    url,
		info:   info,
	}, nil
}

func (hc *httpClient) Harvest(r []Returning) (p []Product, t []Task, e []Error, err error) {
	var b bytes.Buffer
	if err = json.NewEncoder(&b).Encode(Request{
		ClientVersion: 0,
		ServerCookie:  hc.serverCookie,
		MachineInfo:   hc.info,
		Returnings:    r,
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
