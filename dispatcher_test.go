package main

import (
	"testing"

	"github.com/nats-io/gnatsd/server"
	gnatsd "github.com/nats-io/gnatsd/test"

	"github.com/julien/juk/api"
)

func runDefaultServer() *server.Server {
	opts := gnatsd.DefaultTestOptions
	opts.Port = 4666
	return gnatsd.RunServer(&opts)
}

func TestNewDispatcher(t *testing.T) {
	s := runDefaultServer()
	defer s.Shutdown()

	d, err := NewDispatcher("0.0.0.0", 4666)
	if err != nil {
		t.Errorf("got = %v", err)
	}
	defer d.Close()

	if d.conn == nil {
		t.Errorf("got = %v", d.conn)
	}

	d.conn.Publish(api.RegisterJob, &api.JobMsg{
		Name:    "test",
		Path:    "/test",
		Methods: []string{"GET"},
	})

	for {
		select {
		case m := <-d.Messages():
			if m.Name != "test" {
				t.Errorf("got = %v, want test", m.Name)
			}
			return
		}
	}
}

func TestNewDispatcherKO(t *testing.T) {
	d, err := NewDispatcher("0.0.0.0", 4666)
	if err == nil {
		t.Errorf("want an error")
	}
	if d != nil {
		t.Errorf("got = %v, want nil", d)
	}
}

func TestDispatcherRun(t *testing.T) {
	s := runDefaultServer()
	defer s.Shutdown()

	d, err := NewDispatcher("0.0.0.0", 4666)
	if err != nil {
		t.Errorf("got = %v", err)
	}
	defer d.Close()

	if d.conn == nil {
		t.Errorf("got = %v", d.conn)
	}

	d.conn.Publish(api.RegisterJob, &api.JobMsg{
		Name:    "test",
		Path:    "/test",
		Methods: []string{"GET"},
	})

	go d.Run()

	conn, err := api.CreateEncodedConn("0.0.0.0", 4666)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	defer conn.Close()

	jobCh := make(chan *api.Job)
	go func() {
		conn.Subscribe("test", func(m *api.Job) {
			jobCh <- m
		})

		d.Schedule(&api.Job{
			Name: "test",
			Data: []byte("test"),
		})
	}()

	for {
		select {
		case j := <-jobCh:
			if j == nil {
				t.Errorf("want a job")
			}
			if j.Name != "test" {
				t.Errorf(`got = %v, want "test"`, j.Name)
			}
			return
		}
	}
}
