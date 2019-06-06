package server

import (
	"context"
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
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewServer(connectionCount int, host string, port int, handler handleConn) *listening {
	ctx, cancel := context.WithCancel(context.Background())
	return &listening{
		connectionCount: connectionCount,
		h:               handler,
		host:            host,
		port:            port,
		ctx:             ctx,
		cancel:          cancel,
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
		if err := l.h.handle(l.ctx, l.cancel, c); err != nil {
			fmt.Println(err)
		}
	}(conn)
	return err
}
