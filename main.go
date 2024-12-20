package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/julien/juk/api"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	path := DefaultConfigPath
	if len(args) > 1 {
		path = args[1]
	}

	var cfg Config
	if err := cfg.From(path); err != nil {
		fmt.Fprintf(stderr, "%v", err)
		return err
	}

	srv := NewServer(cfg)
	startNatsServer(srv, cfg)
	go startServer(srv, cfg)

	dsp, err := NewDispatcher(cfg.NatsHost, cfg.NatsPort)
	if err != nil {
		fmt.Fprintf(stderr, "%v", err)
		return err
	}
	defer dsp.Close()
	go dsp.Run()

	go handleMessages(dsp, srv)

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintf(stdout, "\nbye\n")
			os.Exit(0)
		}
	}
}

func startServer(s *Server, cfg Config) {
	fmt.Printf("starting HTTP server on %s:%s\n", cfg.Addr, cfg.Port)
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
