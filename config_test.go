package main

import "testing"

func TestLoadConfigOK(t *testing.T) {

	cfg, err := LoadConfig("./config.json")

	if err != nil {
		t.Errorf("got %s\n", err)
	}

	if cfg.Address != DefaultAddr {
		t.Errorf("got %s want %s\n", cfg.Address, DefaultAddr)
	}

	if cfg.Port != DefaultPort {
		t.Errorf("got %s want %s\n", cfg.Port, DefaultPort)
	}

	if cfg.NatsURL != DefaultNatsURL {
		t.Errorf("got %s want %s\n", cfg.NatsURL, DefaultNatsURL)
	}
}

func TestLoadKO(t *testing.T) {
	_, err := LoadConfig("./non-existing.json")
	if err == nil {
		t.Errorf("got %s\n", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Address != DefaultAddr {
		t.Errorf("got %s want %s\n", cfg.Address, DefaultAddr)
	}

	if cfg.Port != DefaultPort {
		t.Errorf("got %s want %s\n", cfg.Port, DefaultPort)
	}

	if cfg.NatsURL != DefaultNatsURL {
		t.Errorf("got %s want %s\n", cfg.NatsURL, DefaultNatsURL)
	}
}
