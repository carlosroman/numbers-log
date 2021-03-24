package pkg

import (
	"fmt"
	"net"
)

type Client struct {
	servAddr string
	conn     *net.TCPConn
}

func NewClient(servAddr string) *Client {
	return &Client{servAddr: servAddr}
}

func (c *Client) Connect() (err error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", c.servAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	c.conn = conn
	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Send(number uint32) error {
	_, err := c.conn.
		Write(
			[]byte(fmt.Sprintf("%09d\n", number)))
	return err
}
