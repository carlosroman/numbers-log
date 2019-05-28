package server

import (
	"golang.org/x/net/netutil"
	"net"
)

type listening struct {
	listener        net.Listener
	connectionCount int
	h               handleConn
}

func newServer(connectionCount int) *listening {
	return &listening{
		connectionCount: connectionCount,
		h:               &handler{},
	}
}

func (l *listening) Start() (err error) {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		return err
	}
	l.listener = netutil.LimitListener(listener, l.connectionCount)
	return err
}

func (l *listening) Stop() (err error) {
	return l.listener.Close()
}

func (l *listening) Process() (err error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return err
	}

	go l.h.handle(conn)
	return err
}
