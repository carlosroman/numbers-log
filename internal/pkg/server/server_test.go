package server

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"sync"
	"testing"
)

func TestStartAndStop(t *testing.T) {
	l := newServer(5)
	defer func() {
		assert.NoError(t, l.Stop())
	}()
	assert.NoError(t, l.Start())
}

func TestStopReturnsError(t *testing.T) {
	l := newServer(1)
	ml := new(mockListener)
	l.listener = ml
	expectedErr := errors.New("some error")
	ml.On("Close").Return(expectedErr)
	assert.EqualError(t, l.Stop(), "some error")
}

func TestProcess(t *testing.T) {
	ml := new(mockListener)
	l := &listening{
		listener: ml,
	}
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
	ml.On("Accept").Return(c, nil)
	assert.NoError(t, l.Process(), "error trying to process connection")
}

type mockListener struct {
	mock.Mock
}

func (m *mockListener) Accept() (net.Conn, error) {
	args := m.Called()
	return args.Get(0).(net.Conn), args.Error(1)
}

func (m *mockListener) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockListener) Addr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}
