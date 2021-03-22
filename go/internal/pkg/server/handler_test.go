package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net"
	"sync"
	"testing"
)

func Test_handler_handle(t *testing.T) {
	logger := getLogger()
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
		name             string
		args             args
		setup            func() (m *mockRepo, h handleConn, l *mockLog)
		write            func(in net.Conn, cxl context.CancelFunc)
		expectConnClosed bool
	}{
		{
			name: "SimpleTest",
			args: a(),
			setup: func() (m *mockRepo, h handleConn, l *mockLog) {
				m = new(mockRepo)
				m.On("IsUnique", uint32(0)).Return(true)
				m.On("IsUnique", uint32(1)).Return(true)
				m.On("IsUnique", uint32(2)).Return(true)

				l = new(mockLog)
				l.On("Info", "000000000", []zapcore.Field(nil))
				l.On("Info", "000000001", []zapcore.Field(nil))
				l.On("Info", "000000002", []zapcore.Field(nil))

				return m, NewHandler(m, l), l
			},
			write: func(in net.Conn, cxl context.CancelFunc) {
				logger.Debug("writing...")
				_, err := in.Write([]byte("000000000\n000000001\n"))
				assert.NoError(t, err, "error writing")
				_, err = in.Write([]byte("000000002\n"))
				assert.NoError(t, err, "error writing")
				logger.Debug("... writing done")
			},
		},
		{
			name: "CtxDone",
			args: a(),
			setup: func() (m *mockRepo, h handleConn, l *mockLog) {
				m = new(mockRepo)
				l = new(mockLog)
				return m, NewHandler(m, l), l
			},
			write: func(in net.Conn, cxl context.CancelFunc) {
				cxl()
				logger.Debug("writing...")
				_, err := in.Write([]byte("000000000\n"))
				assert.EqualError(t, err, io.ErrClosedPipe.Error())
			},
			expectConnClosed: true,
		},
		{
			name: "Terminate",
			args: a(),
			setup: func() (m *mockRepo, h handleConn, l *mockLog) {
				m = new(mockRepo)
				m.On("IsUnique", uint32(0)).Return(true)
				l = new(mockLog)
				l.On("Info", "000000000", []zapcore.Field(nil))
				return m, NewHandler(m, l), l
			},
			write: func(in net.Conn, cxl context.CancelFunc) {
				logger.Debug("writing...")
				_, err := in.Write([]byte("000000000\nterminate\n"))
				assert.NoError(t, err, "error writing")
				_, err = in.Write([]byte("000000001\n"))
				assert.EqualError(t, err, io.ErrClosedPipe.Error())
			},
			expectConnClosed: true,
		},
		{
			name: "DisconnectTooShort",
			args: a(),
			setup: func() (m *mockRepo, h handleConn, l *mockLog) {
				m = new(mockRepo)
				l = new(mockLog)
				return m, NewHandler(m, l), l
			},
			write: func(in net.Conn, cxl context.CancelFunc) {
				logger.Debug("writing...")
				_, err := in.Write([]byte("00000000\n"))
				assert.NoError(t, err, "error writing")
			},
			expectConnClosed: true,
		},
		{
			name: "NotNumber",
			args: a(),
			setup: func() (m *mockRepo, h handleConn, l *mockLog) {
				m = new(mockRepo)
				l = new(mockLog)
				return m, NewHandler(m, l), l
			},
			write: func(in net.Conn, cxl context.CancelFunc) {
				logger.Debug("writing...")
				_, err := in.Write([]byte("ABCDEFGHI\n"))
				assert.NoError(t, err, "error writing")
			},
			expectConnClosed: true,
		},
		{
			name: "DisconnectTooLong",
			args: a(),
			setup: func() (m *mockRepo, h handleConn, l *mockLog) {
				m = new(mockRepo)
				l = new(mockLog)
				return m, NewHandler(m, l), l
			},
			write: func(in net.Conn, cxl context.CancelFunc) {
				logger.Debug("writing...")
				_, err := in.Write([]byte("0000000000\n"))
				assert.NoError(t, err, "error writing")
			},
			expectConnClosed: true,
		},
		{
			name: "Noop",
			args: a(),
			setup: func() (m *mockRepo, h handleConn, l *mockLog) {
				m = new(mockRepo)
				m.On("IsUnique", uint32(0)).Return(false)
				l = new(mockLog)
				return m, NewHandler(m, l), l
			},
			write: func(in net.Conn, cxl context.CancelFunc) {
				logger.Debug("writing...")
				// already present
				_, err := in.Write([]byte("000000000\n"))
				assert.NoError(t, err, "error writing")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				tt.write(tt.args.in, cancel)
			}()

			done := make(chan bool)
			go func() {
				wg.Wait()
				logger.Debug("closing...")
				n, err := tt.args.in.Write([]byte("test"))
				logger.Info("reading conn", zap.Error(err), zap.Int("write", n))
				if tt.expectConnClosed {
					assert.EqualError(t, err, io.ErrClosedPipe.Error())
				} else {
					assert.NoError(t, err)
				}
				assert.NoError(t, tt.args.in.Close())
				assert.NoError(t, tt.args.conn.Close())
				done <- true
			}()
			m, h, l := tt.setup()
			logger.Debug("handling")
			err := h.handle(ctx, cancel, tt.args.conn)
			logger.Debug("done")
			assert.NoError(t, err)
			assert.True(t, <-done)
			m.AssertExpectations(t)
			l.AssertExpectations(t)
		})
	}
}

func getLogger() *zap.Logger {
	logger, err := zap.NewDevelopment(zap.AddCaller())
	if err != nil {
		panic(err)
	}
	return logger
}

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) IsUnique(n uint32) (unique bool) {
	args := m.Called(n)
	return args.Bool(0)
}

func (m *mockRepo) GetReport() string {
	args := m.Called()
	return args.String(0)
}

type mockLog struct {
	mock.Mock
}

func (m *mockLog) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
