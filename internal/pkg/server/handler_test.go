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
		conn, in net.Conn
	}
	a := func() args {
		s, c := net.Pipe()
		return args{
			conn: c,
			in:   s,
		}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "SimpleTest",
			args: a(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println("writing...")
				_, err := tt.args.in.Write([]byte("bob\ncarlos\n"))
				assert.NoError(t, err, "error writing")
				_, err = tt.args.in.Write([]byte("dave\n"))
				assert.NoError(t, err, "error writing")
			}()

			done := make(chan bool)
			go func() {
				wg.Wait()
				assert.NoError(t, tt.args.in.Close())
				assert.NoError(t, tt.args.conn.Close())
				done <- true
			}()

			h := new(handler)
			err := h.handle(tt.args.conn)
			assert.NoError(t, err)
			assert.True(t, <-done)
		})
	}
}
