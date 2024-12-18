package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/julien/juk/api"
)

func main() {
	path := os.Getenv("JUK_CONFIG")
	if strings.TrimSpace(path) == "" {
		if len(os.Args) > 1 {
			path = os.Args[1]
		} else {
			path = "config.json"
		}
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// create a server
	srv := NewServer(cfg)
	startNatsServer(srv, cfg)
	go startServer(srv, cfg)

	// create dispatcher
	dsp, err := NewDispatcher(cfg.NatsHost, cfg.NatsPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	defer dsp.Close()
	go dsp.Run()

	go handleMessages(dsp, srv)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
			fmt.Println("\nstopping")
			os.Exit(0)
		}
	}
}

func startServer(s *Server, cfg Config) {
	fmt.Printf("starting HTTP server on %s:%s\n", cfg.Address, cfg.Port)
	if err := s.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func startNatsServer(s *Server, cfg Config) {
	fmt.Printf("starting NATS server on %s:%d\n", cfg.NatsHost, cfg.NatsPort)
	s.RunNatsServer()
	defer s.ShutdownNatsServer()
}

func createHandler(d *Dispatcher, m *api.JobMsg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			return
		}
		for k, v := range m.ResponseHeaders {
			w.Header().Set(k, v)
		}
		if m.StatusCode != 0 {
			w.WriteHeader(m.StatusCode)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		vars := mux.Vars(r)
		d.Schedule(&api.Job{
			Name:    m.Name,
			Headers: r.Header,
			Data:    body,
			Vars:    vars,
		})
	}
}

func handleMessages(d *Dispatcher, s *Server) {
	for {
		select {
		case m := <-d.Messages():
			s.HandleFunc(m.Path, createHandler(d, m), m.Methods)
		}
	}
}
