package api

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

const (
	RegisterJob = "job.register"
	UpdateJob   = "job.update"
)

type Job struct {
	Data    []byte              `json:"data,omitempty"`
	Name    string              `json:"name"`
	Headers map[string][]string `json:"headers"`
	Vars    map[string]string   `json:"vars"`
}

type JobMsg struct {
	Name            string            `json:"name"`
	Path            string            `json:"path"`
	Methods         []string          `json:"methods,omitempty"`
	ResponseHeaders map[string]string `json:"response_headers,omitempty"`
	StatusCode      int               `json:"status_code,omitempty"`
}

func CreateEncodedConn(host string, port int) (*nats.EncodedConn, error) {

	url := fmt.Sprintf("nats://%s:%d", host, port)

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	conn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
