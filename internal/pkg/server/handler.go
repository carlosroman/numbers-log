package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type handleConn interface {
	handle(conn net.Conn) error
}

type handler struct {
}

func (h *handler) handle(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	for {
		fmt.Println("Reading...")
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
		fmt.Println("msg:'" + strings.TrimRight(msg, "\n") + "'")
		fmt.Println("-----")
	}
}
