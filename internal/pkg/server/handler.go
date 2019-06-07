package server

import (
	"bufio"
	"context"
	"go.uber.org/zap"
	"io"
	"net"
	"strconv"
	"strings"
)

type log interface {
	Info(msg string, fields ...zap.Field)
}

type handleConn interface {
	handle(ctx context.Context, cancel context.CancelFunc, conn net.Conn) error
}

type handler struct {
	nc     NumberChecker
	logger log
}

func NewHandler(numberChecker NumberChecker, logger log) handleConn {
	return &handler{
		nc:     numberChecker,
		logger: logger,
	}
}

func (h *handler) handle(ctx context.Context, cancel context.CancelFunc, conn net.Conn) error {
	reader := bufio.NewReader(conn)
	c := make(chan string, 1)
	e := make(chan error, 1)
	d := make(chan struct{}, 1)

	for {
		go func() {
			msg, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF || err == io.ErrClosedPipe {
					d <- struct{}{}
				}
				e <- err
			}
			c <- msg
		}()
		select {
		case <-ctx.Done():
			if errConn := conn.Close(); errConn != nil {
				return errConn
			}
			return nil
		case <-d:
			return nil
		case err := <-e:
			if errConn := conn.Close(); errConn != nil {
				return errConn
			}
			return err
		case msg := <-c:
			v := strings.TrimRight(msg, "\n")
			if v == "terminate" {
				cancel()
				continue
			}
			if len(v) != 9 {
				if errConn := conn.Close(); errConn != nil {
					// log errConn
					return errConn
				}
				return nil
			}

			i, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				if errConn := conn.Close(); errConn != nil {
					// log errConn
					return errConn
				}
				return nil
			}

			if h.nc.IsUnique(uint32(i)) {
				h.logger.Info(v)
			}
		}
	}
}
