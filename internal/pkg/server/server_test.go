package server

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
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
	ctx, cancel := context.WithCancel(context.Background())
	ml := new(mockListener)
	hm := new(mockHandleConn)
	l := &listening{
		listener: ml,
		h:        hm,
		ctx:      ctx,
		cancel:   cancel,
	}
	s, c := net.Pipe()
	defer func() {
		assert.NoError(t, s.Close())
		assert.NoError(t, c.Close())
	}()
	ml.On("Accept").Return(c, nil)
	hm.On("handle", ctx, mock.AnythingOfType("context.CancelFunc"), c).Return(nil).
		Run(func(args mock.Arguments) {
			cancel()
		})
	err := l.Process()
	assert.NoError(t, err, "error trying to process connection")
	ml.AssertExpectations(t)
	hm.AssertExpectations(t)
}
func TestProcess_handle_error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ml := new(mockListener)
	hm := new(mockHandleConn)
	l := &listening{
		listener: ml,
		h:        hm,
		ctx:      ctx,
		cancel:   cancel,
	}
	s, c := net.Pipe()
	defer func() {
		assert.NoError(t, s.Close())
		assert.NoError(t, c.Close())
	}()
	ml.On("Accept").Return(c, nil)
	expErr := errors.New("some error")
	hm.On("handle", ctx, mock.AnythingOfType("context.CancelFunc"), c).Return(expErr)
	err := l.Process()
	assert.Errorf(t, err, expErr.Error(), "expected an error")
	ml.AssertExpectations(t)
	hm.AssertExpectations(t)
}

func TestProcess_accept_error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ml := new(mockListener)
	hm := new(mockHandleConn)
	l := &listening{
		listener: ml,
		h:        hm,
		ctx:      ctx,
		cancel:   cancel,
	}
	s, c := net.Pipe()
	defer func() {
		assert.NoError(t, s.Close())
		assert.NoError(t, c.Close())
	}()
	expErr := errors.New("some error")
	ml.On("Accept").Return(c, expErr)
	err := l.Process()
	assert.Errorf(t, err, expErr.Error(), "expected an error")
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
