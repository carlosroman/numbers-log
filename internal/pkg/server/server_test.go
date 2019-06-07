package server

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"sync"
	"testing"
)

func TestStartAndStop(t *testing.T) {
	l := NewServer(5, "127.0.0.1", 0, &handler{})
	defer func() {
		assert.NoError(t, l.Stop())
	}()
	assert.NoError(t, l.Start())
}

func TestStopReturnsError(t *testing.T) {
	l := NewServer(1, "127.0.0.1", 0, &handler{})
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
	var once sync.Once
	ml.On("Accept").Return(getConn(), nil)
	hm.On("handle", mock.Anything, mock.AnythingOfType("context.CancelFunc"), mock.Anything).Return(nil).
		Run(func(args mock.Arguments) {
			once.Do(args.Get(1).(context.CancelFunc))
		})
	err := l.Process()
	assert.NoError(t, err, "error trying to process connection")
	ml.AssertExpectations(t)
	hm.AssertExpectations(t)
}

func TestProcess_handle_error(t *testing.T) {
	ml := new(mockListener)
	hm := new(mockHandleConn)
	l := &listening{
		listener: ml,
		h:        hm,
	}
	ml.On("Accept").Return(getConn(), nil)
	hm.On("handle", mock.Anything, mock.AnythingOfType("context.CancelFunc"), mock.Anything).Return(errors.New("some error"))
	err := l.Process()
	assert.Errorf(t, err, "some error", "expected an error")
	ml.AssertExpectations(t)
	hm.AssertExpectations(t)
}

func getConn() net.Conn {
	server, client := net.Pipe()
	defer server.Close()
	return client
}

func TestProcess_accept_error(t *testing.T) {
	ml := new(mockListener)
	hm := new(mockHandleConn)
	l := &listening{
		listener: ml,
		h:        hm,
	}
	ml.On("Accept").Return(getConn(), errors.New("some error"))
	err := l.Process()
	assert.Errorf(t, err, "some error", "expected an error")
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

func (m *mockHandleConn) handle(ctx context.Context, cancel context.CancelFunc, conn net.Conn) error {
	return m.Called(ctx, cancel, conn).Error(0)
}
