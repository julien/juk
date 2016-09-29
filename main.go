package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

var (
	envConfig = os.Getenv("JUK_CONFIG")
)

func main() {

	var path string
	if envConfig != "" {
		path = envConfig
	} else if len(os.Args) == 2 {
		path = os.Args[1]
	} else {
		path = "config.json"
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// create dispatcher
	dsp, err := NewDispatcher(cfg.NatsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	defer dsp.Close()
	go dsp.Run()

	// create a server
	srv := NewServer(cfg)
	go startServer(srv)

	go handleMessages(dsp, srv)

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
			fmt.Println("\nstopping")
			os.Exit(0)
		}
	}

}

func startServer(s *Server) {
	fmt.Println("starting server")
	if err := s.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func createHandler(d *Dispatcher, m *JobMsg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
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
		d.Schedule(&Job{
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
