package main

import (
	"fmt"
	"log"

	"github.com/julien/juk/api"
)

const exampleJob = "example.job"

func main() {
	conn, err := api.CreateEncodedConn("localhost", 4222)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	msg := api.JobMsg{
		Name:    exampleJob,
		Path:    "/example/{name}/{id:[0-9]+}",
		Methods: []string{"GET", "POST"},
		ResponseHeaders: map[string]string{
			"x-token":     "fwaa5ckoCWDOIawjd",
			"x-signature": "aksjakjSAHQGDBC-123",
		},
		StatusCode: 201,
	}
	conn.Publish(api.RegisterJob, msg)
	conn.Subscribe(exampleJob, execJob)

	done := make(chan bool)
	<-done
}

func execJob(m *api.Job) {
	fmt.Printf("got a job\n%+v\n", m)
	fmt.Println(m.Data)
	// Do what you want with the job here
}
