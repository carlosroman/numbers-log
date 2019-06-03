package server

import (
	"bufio"
	"go.uber.org/zap"
	"io"
	"net"
	"strconv"
	"strings"
)

type add interface {
	Add(n uint32) (unique bool)
}

type log interface {
	Info(msg string, fields ...zap.Field)
}

type handleConn interface {
	handle(conn net.Conn) error
}

type handler struct {
	repo   add
	logger log
}

func NewHandler(repo add, logger log) handleConn {
	return &handler{
		repo:   repo,
		logger: logger,
	}
}

func (h *handler) handle(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	for {
		//if err = conn.SetReadDeadline(time.Now().Add(1* time.Second)); err!=nil{
		//	return err
		//}

		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF || err == io.ErrClosedPipe {
				return nil
			}
			conn.Close()
			return err
		}
		v := strings.TrimRight(msg, "\n")
		if len(v) != 9 {
			continue
		}

		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			continue
		}
		if h.repo.Add(uint32(i)) {
			h.logger.Info(v)
		}
	}
}