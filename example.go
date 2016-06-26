package main

import (
	"fmt"
	"log"
)

const exampleJob = "example.job"

func Example() {
	conn, err := CreateEncodedConn("0.0.0.0:4222")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	msg := JobMsg{
		Name:    exampleJob,
		Path:    "/example/{name}/{id:[0-9]+}",
		Methods: []string{"GET", "POST"},
		ResponseHeaders: map[string]string{
			"x-token":     "fwaa5ckoCWDOIawjd",
			"x-signature": "aksjakjSAHQGDBC-123",
		},
		StatusCode: 201,
	}
	conn.Publish(RegisterJob, msg)
	conn.Subscribe(exampleJob, execJob)

	done := make(chan bool)
	<-done
}

func execJob(m *Job) {
	fmt.Printf("got a job\n%+v\n", m)
	// Do what you want with the job here
}
