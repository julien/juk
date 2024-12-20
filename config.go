package main

import (
	"encoding/json"
	"os"
)

const (
	DefaultAddr       = "127.0.0.1"
	DefaultConfigPath = "config.json"
	DefaultPort       = "8000"
	DefaultNatsHost   = "127.0.0.1"
	DefaultNatsPort   = 4222
	DefaultLogLevel   = "DEBUG"
)

type Config struct {
	Addr  string `json:"addr"`
	Port     string `json:"port"`
	Secure   bool   `json:"secure,omitempty"`
	Certfile string `json:"certfile,omitempty"`
	Keyfile  string `json:"keyfile,omitempty"`
	NatsHost string `json:"nats_host"`
	NatsPort int    `json:"nats_port"`
	LogLevel string `json:"log_level"`
}

func (c *Config) From(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(c)
}
