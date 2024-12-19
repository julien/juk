package api

import (
	"testing"

	gnatsd "github.com/nats-io/gnatsd/test"
)

func TestCreateEncodedConn(t *testing.T) {
	opts := gnatsd.DefaultTestOptions
	opts.Port = 6666
	s := gnatsd.RunServer(&opts)
	defer s.Shutdown()

	c, err := CreateEncodedConn("localhost", 6666)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	defer c.Close()

}
