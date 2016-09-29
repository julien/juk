package main

import (
	"net/http"
	"sync"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	"github.com/gorilla/mux"
)

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Methods []string
}

type Server struct {
	secure   bool
	certFile string
	keyFile  string
	mu       sync.Mutex
	routes   map[string]Route
	server   *graceful.Server
}

func NewServer(cfg *Config) *Server {
	return &Server{
		secure:   cfg.Secure,
		certFile: cfg.Certfile,
		keyFile:  cfg.Keyfile,
		routes:   make(map[string]Route),
		server: &graceful.Server{
			Timeout: 10 * time.Second,
			Server: &http.Server{
				Addr:    cfg.Address + ":" + cfg.Port,
				Handler: mux.NewRouter(),
			},
		},
	}
}

func (s *Server) Start() error {
	if s.secure && s.certFile != "" && s.keyFile != "" {
		return s.server.ListenAndServeTLS(s.certFile, s.keyFile)
	} else {
		return s.server.ListenAndServe()
	}
}

func (s *Server) HandleFunc(path string, handler http.HandlerFunc, methods []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.routes[path]; ok {
		return
	}
	s.routes[path] = Route{path, handler, methods}

	// Create a mux route
	r := s.server.Server.Handler.(*mux.Router)
	r.HandleFunc(path, handler)
	r.Methods(methods...)
}
