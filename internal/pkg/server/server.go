package server

import (
	"fmt"
	"golang.org/x/net/netutil"
	"net"
)

type listening struct {
	listener        net.Listener
	connectionCount int
	h               handleConn
	host            string
	port            int
}

func NewServer(connectionCount int, host string, port int, handler handleConn) *listening {
	return &listening{
		connectionCount: connectionCount,
		h:               handler,
		host:            host,
		port:            port,
	}
}

func (l *listening) Start() (err error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%v", l.host, l.port))
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
		fmt.Println(err)
		return err
	}

	go func(c net.Conn) {
		if err := l.h.handle(c); err != nil {
			fmt.Println(err)
		}
	}(conn)
	return err
}
