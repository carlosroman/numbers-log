package server

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"sync"
	"testing"
)

func Test_handler_handle(t *testing.T) {
	logger := getLogger()
	defer logger.Sync()
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
		name  string
		args  args
		setup func() (m *mockRepo, h handleConn, l *mockLog)
		write func(in net.Conn)
	}{
		{
			name: "SimpleTest",
			args: a(),
			setup: func() (m *mockRepo, h handleConn, l *mockLog) {
				m = new(mockRepo)
				m.On("Add", uint32(0)).Return(false)
				m.On("Add", uint32(1)).Return(false)
				m.On("Add", uint32(2)).Return(false)

				l = new(mockLog)
				l.On("Info", "000000000", []zapcore.Field(nil))
				l.On("Info", "000000001", []zapcore.Field(nil))
				l.On("Info", "000000002", []zapcore.Field(nil))

				return m, newHandler(m, l), l
			},
			write: func(in net.Conn) {
				fmt.Println("writing...")
				_, err := in.Write([]byte("000000000\n000000001\n"))
				assert.NoError(t, err, "error writing")
				_, err = in.Write([]byte("000000002\n"))
				assert.NoError(t, err, "error writing")
			},
		},
		{
			name: "Noop",
			args: a(),
			setup: func() (m *mockRepo, h handleConn, l *mockLog) {
				m = new(mockRepo)
				m.On("Add", uint32(0)).Return(true)

				l = new(mockLog)
				l.On("Info", "000000000", []zapcore.Field(nil))

				return m, newHandler(m, l), l
			},
			write: func(in net.Conn) {
				fmt.Println("writing...")
				_, err := in.Write([]byte("000000000\n"))
				assert.NoError(t, err, "error writing")
				_, err = in.Write([]byte("0000000001\n"))
				assert.NoError(t, err, "error writing")
				_, err = in.Write([]byte("00000002\n"))
				assert.NoError(t, err, "error writing")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				tt.write(tt.args.in)
			}()

			done := make(chan bool)
			go func() {
				wg.Wait()
				assert.NoError(t, tt.args.in.Close())
				assert.NoError(t, tt.args.conn.Close())
				done <- true
			}()

			m, h, l := tt.setup()
			err := h.handle(tt.args.conn)
			assert.NoError(t, err)
			assert.True(t, <-done)
			m.AssertExpectations(t)
			l.AssertExpectations(t)
		})
	}
}

func getLogger() *zap.Logger {
	rawJSON := []byte(`{
	  "level": "info",
	  "encoding": "console",
	  "outputPaths": ["stdout"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Add(n uint32) (unique bool) {
	args := m.Called(n)
	return args.Bool(0)
}

type mockLog struct {
	mock.Mock
}

func (m *mockLog) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
