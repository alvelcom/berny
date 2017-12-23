package api

import "encoding/json"

type Request struct {
	ClientVersion int    `json:"client_version"`
	ServerCookie  string `json:"server_cookie,omitempty"`

	MachineInfo MachineInfo `json:"machine_info"`
	Returnings  []Returning `json:"returninds"`
}

type MachineInfo struct {
	IPs  []string `json:"ips"`
	FQDN string   `json:"fqdn"`

	// Optional fields, useful for templates and verification
	Host   string `json:"host,omitempty"`
	Domain string `json:"domain,omitempty"`
	// Host+Domain == FQDN

	Cluster  string `json:"cluster,omitempty"`
	NodeType string `json:"node_type,omitempty"`
	Id       string `json:"id,omitempty"`
	// Usually Cluster+NodeType+Id == Host

	// Cloud-specific
	Provider string `json:"provider,omitempty"` // aws, gcp or your own rasberry under the desk
	Region   string `json:"region,omitempty"`

	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	Geo     string `json:"geo,omitempty"` // free format

	Extra map[string]string `json:"extra,omitempty"`
}

type Returning struct {
	Name []string
	Type string
	Body json.RawMessage
}

type Response struct {
	ServerVersion int    `json:"server_version"`
	ServerCookie  string `json:"server_cookie"`

	Errors   []Error   `json:"errors"`
	Tasks    []Task    `json:"tasks"`
	Products []Product `json:"products"`
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Task struct {
	Name []string        `json:"name"`
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type Product struct {
	Name []string `json:"name"`
	Mask int      `json:"mask"`
	Body string   `json:"body"`
}

func (ci *MachineInfo) Verify() error {
	return nil
}
