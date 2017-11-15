package api

import "encoding/json"

type Request struct {
	ClientVersion int    `json:"client-version"`
	ServerCookie  string `json:"server-cookie,omitempty"`

	MachineInfo MachineInfo `json:"machine-info"`
	Returnings  []Returning
}

type MachineInfo struct {
	IPs  []string
	FQDN string

	// Optional fields, useful for templates and verification
	Host   string `json:",omitempty"`
	Domain string `json:",omitempty"`
	// Host+Domain == FQDN

	Cluster  string `json:",omitempty"`
	NodeType string `json:"node-type,omitempty"`
	Id       string `json:",omitempty"`
	// Usually Cluster+NodeType+Id == Host

	// Cloud-specific
	Provider string `json:",omitempty"` // aws, gcp or your own rasberry under the desk
	Region   string `json:",omitempty"`

	City    string `json:",omitempty"`
	Country string `json:",omitempty"`
	Geo     string `json:",omitempty"` // free format

	Extra map[string]string `json:",omitempty"`
}

type Returning struct {
	Name []string
	Type string
	Body json.RawMessage
}

type Response struct {
	ServerVersion int    `json:"server-version"`
	ServerCookie  string `json:"server-cookie"`

	Errors   []Error
	Tasks    []Task
	Products []Product
}

type Error struct {
	Type    string
	Message string
}

type Task struct {
	Name []string
	Type string
	Body json.RawMessage
}

type Product struct {
	Name []string
	Mask int
	Body string
}

func (ci *MachineInfo) Verify() error {
	return nil
}
