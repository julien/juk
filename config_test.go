package main

import "testing"

func TestLoadConfigOK(t *testing.T) {
	var cfg Config
	if err := cfg.From("./config.json"); err != nil {
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

	if cfg.LogLevel != DefaultLogLevel {
		t.Errorf("got = %v, want = %v", cfg.LogLevel, DefaultLogLevel)
	}
}

func TestLoadConfigKO(t *testing.T) {
	var cfg Config
	if err := cfg.From("./non-existing.json"); err == nil {
		t.Errorf("want an error")
	}
}
