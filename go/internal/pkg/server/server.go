package server

import (
	"context"
	"fmt"
	"golang.org/x/net/netutil"
	"net"
	"time"
)

type listening struct {
	listener        net.Listener
	connectionCount int
	h               handleConn
	host            string
	port            int
	tickerDuration  time.Duration
}

func NewServer(connectionCount int, host string, port int, handler handleConn, tickerDuration time.Duration) *listening {
	return &listening{
		connectionCount: connectionCount,
		h:               handler,
		host:            host,
		port:            port,
		tickerDuration:  tickerDuration,
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
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan net.Conn, 1)
	e := make(chan error, 1)
	ticker := time.NewTicker(l.tickerDuration)

	go func() {
		for {
			select {
			case <-ticker.C:
				l.h.printReport()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	for {
		go func() {
			conn, err := l.listener.Accept()
			if err != nil {
				e <- err
			} else {
				c <- conn
			}
		}()
		select {
		case <-ctx.Done():
			return nil
		case conn := <-c:
			go func(cn net.Conn) {
				if err := l.h.handle(ctx, cancel, cn); err != nil {
					e <- err
				}
			}(conn)
		case err := <-e:
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
}
