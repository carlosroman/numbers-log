package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
)

func Test_handler_handle(t *testing.T) {
	type args struct {
		conn net.Conn
	}
	tests := []struct {
		name string
		h    *handler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{}
			h.handle(tt.args.conn)
		})
	}
}

func TestProcessHandler(t *testing.T) {
	h := new(handler)
	s, c := net.Pipe()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("writing...")
		_, err := s.Write([]byte("bob\ncarlos\n"))
		assert.NoError(t, err, "error writing")
		_, err = s.Write([]byte("dave\n"))
		assert.NoError(t, err, "error writing")
	}()
	go func() {
		wg.Wait()
		assert.NoError(t, s.Close())
		assert.NoError(t, c.Close())
	}()
	fmt.Println("reading...")
	assert.NoError(t, h.handle(c), "error trying to process connection")
}
