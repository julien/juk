package main

import (
	"testing"

	"github.com/nats-io/gnatsd/server"
	gnatsd "github.com/nats-io/gnatsd/test"
)

func runDefaultServer() *server.Server {
	opts := gnatsd.DefaultTestOptions
	opts.Port = 4666
	return gnatsd.RunServer(&opts)
}

func TestNewDispatcher(t *testing.T) {
	s := runDefaultServer()
	defer s.Shutdown()

	d, err := NewDispatcher("0.0.0.0:4666")
	if err != nil {
		t.Errorf("got %v want a dispatcher\n", err)
	}
	defer d.Close()

	if d.conn == nil {
		t.Errorf("got %v, want a connection\n", d.conn)
	}

	d.conn.Publish(RegisterJob, &JobMsg{
		Name:    "test",
		Path:    "/test",
		Methods: []string{"GET"},
	})

	for {
		select {
		case m := <-d.Messages():
			if m.Name != "test" {
				t.Errorf("got %s want test\n", m.Name)
			}
			return
		}
	}

}

func TestNewDispatcherKO(t *testing.T) {

	d, err := NewDispatcher("0.0.0.0:4666")
	if err == nil {
		t.Errorf("want an error\n")
	}
	if d != nil {
		t.Errorf("want dispatcher to be nil\n")
	}

}

func TestDispatcherRun(t *testing.T) {
	s := runDefaultServer()
	defer s.Shutdown()

	d, err := NewDispatcher("0.0.0.0:4666")
	if err != nil {
		t.Errorf("got %v want a dispatcher\n", err)
	}
	defer d.Close()

	if d.conn == nil {
		t.Errorf("got %v, want a connection\n", d.conn)
	}

	d.conn.Publish(RegisterJob, &JobMsg{
		Name:    "test",
		Path:    "/test",
		Methods: []string{"GET"},
	})

	go d.Run()

	conn, err := CreateEncodedConn("0.0.0.0:4666")
	if err != nil {
		t.Errorf("%s\n", err)
	}
	defer conn.Close()

	jobCh := make(chan *Job)
	go func() {
		conn.Subscribe("test", func(m *Job) {
			jobCh <- m
		})

		d.Schedule(&Job{
			Name: "test",
			Data: []byte("test"),
		})
	}()

	for {
		select {
		case j := <-jobCh:
			if j == nil {
				t.Errorf("want a job\n")
			}

			if j.Name != "test" {
				t.Errorf("want \"test\", got %s\n", j.Name)
			}

			return
		}
	}

}
