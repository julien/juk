package main

import (
	"encoding/json"
	"os"
)

const (
	DefaultAddr     = "127.0.0.1"
	DefaultPort     = "8000"
	DefaultNatsHost = "127.0.0.1"
	DefaultNatsPort = 4222
)

// Config holds the necessary values needed to start the HTTP server
// connect to the NATS server and store the information.
type Config struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Secure   bool   `json:"secure,omitempty"`
	Certfile string `json:"certfile,omitempty"`
	Keyfile  string `json:"keyfile,omitempty"`
	NatsHost string `json:"nats_host"`
	NatsPort int    `json:"nats_port"`
}

// LoadConfig returns a Config instance from a given file name or
// an error if the file could not be read or decoded.
func LoadConfig(name string) (Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	cfg := Config{}
	dec := json.NewDecoder(f)
	if err = dec.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// DefaultConfig returns a Config instance with default values.
func DefaultConfig() Config {
	return Config{
		Address:  DefaultAddr,
		Port:     DefaultPort,
		NatsHost: DefaultNatsHost,
		NatsPort: DefaultNatsPort,
	}
}
