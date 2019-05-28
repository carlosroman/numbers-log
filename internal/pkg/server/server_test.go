package server

import (
	"errors"
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
	hm := new(mockHandleConn)
	l := &listening{
		listener: ml,
		h:        hm,
	}
	s, c := net.Pipe()
	defer func() {
		assert.NoError(t, s.Close())
		assert.NoError(t, c.Close())
	}()
	ml.On("Accept").Return(c, nil).Once()
	wg := sync.WaitGroup{}
	wg.Add(1)
	hm.On("handle", c).Return(nil).Run(func(args mock.Arguments) {
		wg.Done()
	})
	assert.NoError(t, l.Process(), "error trying to process connection")
	wg.Wait()
	ml.AssertExpectations(t)
	hm.AssertExpectations(t)
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

type mockHandleConn struct {
	mock.Mock
}

func (m *mockHandleConn) handle(conn net.Conn) error {
	return m.Called(conn).Error(0)
}
