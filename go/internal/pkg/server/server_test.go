package server

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"sync"
	"testing"
	"time"
)

func TestStartAndStop(t *testing.T) {
	l := NewServer(5, "127.0.0.1", 0, &handler{}, time.Minute)
	defer func() {
		assert.NoError(t, l.Stop())
	}()
	assert.NoError(t, l.Start())
}

func TestStopReturnsError(t *testing.T) {
	l := NewServer(1, "127.0.0.1", 0, &handler{}, time.Minute)
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
		listener:       ml,
		h:              hm,
		tickerDuration: time.Minute,
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
		listener:       ml,
		h:              hm,
		tickerDuration: time.Minute,
	}
	ml.On("Accept").Return(getConn(), nil)
	hm.On("handle", mock.Anything, mock.AnythingOfType("context.CancelFunc"), mock.Anything).Return(errors.New("some error"))
	err := l.Process()
	assert.Errorf(t, err, "some error", "expected an error")
	ml.AssertExpectations(t)
	hm.AssertExpectations(t)
}

func TestProcess_accept_error(t *testing.T) {
	ml := new(mockListener)
	hm := new(mockHandleConn)
	l := &listening{
		listener:       ml,
		h:              hm,
		tickerDuration: time.Minute,
	}
	ml.On("Accept").Return(getConn(), errors.New("some error"))
	err := l.Process()
	assert.Errorf(t, err, "some error", "expected an error")
	ml.AssertExpectations(t)
	hm.AssertExpectations(t)
}

func TestProcess_tick(t *testing.T) {
	ml := new(mockListener)
	hm := new(mockHandleConn)
	l := &listening{
		listener:       ml,
		h:              hm,
		tickerDuration: time.Second,
	}

	asertWg := sync.WaitGroup{}
	asertWg.Add(1)

	ml.On("Accept").Return(getConn(), nil)
	var printOne = sync.Once{}
	hm.On("printReport").Return().Run(func(args mock.Arguments) {
		defer printOne.Do(asertWg.Done)
	})

	hmWg := sync.WaitGroup{}
	hmWg.Add(1)
	hm.On("handle", mock.Anything, mock.AnythingOfType("context.CancelFunc"), mock.Anything).Return(errors.New("some error")).Run(func(args mock.Arguments) {
		hmWg.Wait()
	})

	go l.Process()

	asertWg.Wait()
	ml.AssertExpectations(t)
	hmWg.Done()
}

func getConn() net.Conn {
	server, client := net.Pipe()
	defer server.Close()
	return client
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

func (m *mockHandleConn) printReport() {
	m.Called()
}

func (m *mockHandleConn) handle(ctx context.Context, cancel context.CancelFunc, conn net.Conn) error {
	return m.Called(ctx, cancel, conn).Error(0)
}
