package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				_, err := tt.args.in.Write([]byte("000000000\n000000001\n"))
				assert.NoError(t, err, "error writing")
				_, err = tt.args.in.Write([]byte("000000002\n"))
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
			mrepo := new(mockRepo)
			h.repo = mrepo
			mrepo.On("Add", uint32(0)).Return(false)
			mrepo.On("Add", uint32(1)).Return(false)
			mrepo.On("Add", uint32(2)).Return(false)
			err := h.handle(tt.args.conn)
			assert.NoError(t, err)
			assert.True(t, <-done)
			mrepo.AssertExpectations(t)
		})
	}
}

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Add(n uint32) (unique bool) {
	args := m.Called(n)
	return args.Bool(0)
}
