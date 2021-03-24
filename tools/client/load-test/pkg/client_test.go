package pkg_test

import (
	"bufio"
	"bytes"
	"load-test/pkg"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Connect(t *testing.T) {
	l, err := startTestServer()
	require.NoError(t, err)
	defer func() { _ = l.Close() }()
	client := pkg.NewClient(l.Addr().String())
	require.NoError(t, client.Connect())
	defer func() { _ = client.Close() }()
}

func TestClient_Connect_err(t *testing.T) {
	client := pkg.NewClient("0.0.0.0:999999")
	assert.Error(t, client.Connect())
}

func TestClient_Send(t *testing.T) {
	l, err := startTestServer()
	require.NoError(t, err)
	defer func() { _ = l.Close() }()
	resp := make(chan string)
	go func() {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		require.NoError(t, err)
		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		require.NoError(t, err)
		buffer := bufio.NewReader(bytes.NewBuffer(buf))
		line,_, err := buffer.ReadLine()
		require.NoError(t, err)
		resp <- string(line)
	}()

	client := pkg.NewClient(l.Addr().String())
	require.NoError(t, client.Connect())
	defer func() { _ = client.Close() }()

	//007007009
	err = client.Send(7007009)
	require.NoError(t, err)
	res := <-resp
	assert.Equal(t, "007007009", res)
}

func startTestServer() (net.Listener, error) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	u, err := url.Parse(s.URL)
	if err != nil {
		return nil, err
	}
	s.Close()
	l, err := net.Listen("tcp", u.Host)
	return l, err
}
