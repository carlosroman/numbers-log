package server

import (
	"bufio"
	"fmt"
	"golang.org/x/net/netutil"
	"io"
	"net"
	"strings"
)

type listening struct {
	listener        net.Listener
	connectionCount int
}

func newServer(connectionCount int) *listening {
	return &listening{
		connectionCount: connectionCount,
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

	go func(c net.Conn) {
		reader := bufio.NewReader(c)
		for {
			fmt.Println("Reading...")
			//if err = conn.SetReadDeadline(time.Now().Add(1* time.Second)); err!=nil{
			//	return err
			//}

			msg, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF || err == io.ErrClosedPipe {
					return
				}
				c.Close()
				return
			}
			fmt.Println("msg:'" + strings.TrimRight(msg, "\n") + "'")
			fmt.Println("-----")
		}
	}(conn)
	return err
}
