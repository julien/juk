package main

import (
	"fmt"

	"github.com/julien/juk/api"
	"github.com/nats-io/nats"
)

// Dispatcher is responsable for registering new Job hooks and
// dispatching Jobs to the corresponding channels over NATS.
type Dispatcher struct {
	channels map[string]chan *api.Job
	conn     *nats.EncodedConn
	jobCh    chan *api.Job
	messages chan *api.JobMsg
}

// NewDispatcher returns a new Dispatcher instance given a URL
// or an error if the connection to NATS fails.
func NewDispatcher(host string, port int) (*Dispatcher, error) {

	conn, err := api.CreateEncodedConn(host, port)
	if err != nil {
		fmt.Println("connection error, oh shit: %s\n", err)
		return nil, err
	}

	dsp := &Dispatcher{
		channels: make(map[string]chan *api.Job),
		conn:     conn,
		jobCh:    make(chan *api.Job),
		messages: make(chan *api.JobMsg),
	}

	dsp.conn.Subscribe(api.RegisterJob, func(m *api.JobMsg) {

		if _, ok := dsp.channels[m.Name]; ok {
			return
		}

		dsp.channels[m.Name] = make(chan *api.Job)
		dsp.conn.BindSendChan(m.Name, dsp.channels[m.Name])

		// Notify
		dsp.messages <- m
	})

	return dsp, nil
}

func (d *Dispatcher) Close() {
	if d.conn != nil {
		d.conn.Close()
	}
}

func (d *Dispatcher) Messages() <-chan *api.JobMsg {
	return d.messages
}

// Schedule sends a Job to the Dispatcher's job channel.
func (d *Dispatcher) Schedule(j *api.Job) {
	d.jobCh <- j
}

// Run periodically checks the Dispatcher's job channel
// and redirects incoming Jobs to the correponding channel.
func (d *Dispatcher) Run() {
	for {
		select {
		case j := <-d.jobCh:
			go func(j *api.Job) {
				ch, ok := d.channels[j.Name]
				if ok {
					ch <- j
				}
			}(j)
		}
	}
}
