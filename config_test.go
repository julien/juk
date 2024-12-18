package main

import "testing"

func TestLoadConfigOK(t *testing.T) {
	cfg, err := LoadConfig("./config.json")
	if err != nil {
		t.Errorf("got = %v", err)
	}

	if cfg.Address != DefaultAddr {
		t.Errorf("got = %v, want = %v", cfg.Address, DefaultAddr)
	}

	if cfg.Port != DefaultPort {
		t.Errorf("got = %v, want = %v", cfg.Port, DefaultPort)
	}

	if cfg.NatsHost != DefaultNatsHost {
		t.Errorf("got = %v, want = %v", cfg.NatsHost, DefaultNatsHost)
	}

	if cfg.NatsPort != DefaultNatsPort {
		t.Errorf("got = %v, want = %v", cfg.NatsPort, DefaultNatsPort)
	}
}

func TestLoadKO(t *testing.T) {
	if _, err := LoadConfig("./non-existing.json"); err == nil {
		t.Errorf("want an error")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Address != DefaultAddr {
		t.Errorf("got = %v, want = %v", cfg.Address, DefaultAddr)
	}

	if cfg.Port != DefaultPort {
		t.Errorf("got = %v, want = %v", cfg.Port, DefaultPort)
	}

	if cfg.NatsHost != DefaultNatsHost {
		t.Errorf("got = %v, want = %v", cfg.NatsHost, DefaultNatsHost)
	}

	if cfg.NatsPort != DefaultNatsPort {
		t.Errorf("got = %v, want = %v", cfg.NatsPort, DefaultNatsPort)
	}
}
