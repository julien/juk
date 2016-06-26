package main

import (
	"net/http"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	c := &Config{
		Address: DefaultAddr,
		Port:    "3333",
	}

	s := NewServer(c)
	if s == nil {
		t.Errorf("want a server, got nil")
	}
}

func TestServerStart(t *testing.T) {
	c := &Config{
		Address: DefaultAddr,
		Port:    "3333",
	}

	s := NewServer(c)
	if s == nil {
		t.Errorf("want a server, got nil")
	}

	errCh := make(chan error)

	go func(errCh chan error) {
		errCh <- s.Start()
	}(errCh)

	timeoutCh := time.After(1 * time.Second)
	quitCh := time.After(2 * time.Second)

	for {
		select {
		case e := <-errCh:
			t.Errorf("%s\n", e)

		case <-timeoutCh:
			s.server.Timeout = 0

		case <-quitCh:
			return
		}
	}
}

func TestServerHandleFunc(t *testing.T) {
	c := &Config{
		Address: DefaultAddr,
		Port:    "4444",
	}

	s := NewServer(c)
	if s == nil {
		t.Errorf("want a server, got nil")
	}

	errCh := make(chan error)

	go func(errCh chan error) {
		errCh <- s.Start()
	}(errCh)

	handleFoo := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}

	handlerCh := time.After(1 * time.Second)
	reqCh := time.After(2 * time.Second)
	quitCh := time.After(4 * time.Second)

	for {
		select {
		case e := <-errCh:
			t.Errorf("%s\n", e)

		case <-handlerCh:

			s.HandleFunc("/foo", handleFoo, []string{"GET"})

		case <-reqCh:

			req, err := http.NewRequest("GET",
				"http://localhost:4444/foo", nil)
			if err != nil {
				t.Errorf("%s\n", err)
			}

			client := http.DefaultClient
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("%s\n", err)
			}

			if resp.StatusCode != 200 {
				t.Errorf("got %d want 200\n", resp.StatusCode)
			}

			s.server.Timeout = 0

		case <-quitCh:
			return
		}
	}

}
