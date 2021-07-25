package postern

import (
	"net"
)

type DarkConn struct {
	net.Conn
	*Cipher
}

func NewConn(conn net.Conn, cipher *Cipher) *DarkConn {
	return &DarkConn{conn, cipher}
}

func DarkDial(dest, remote string, cipher *Cipher) (conn *DarkConn, err error) {
	c, err := net.Dial("tcp", remote)
	if err != nil {
		return
	}
	conn = NewConn(c, cipher)
	if _, err := conn.Write(GenHead(dest)); err != nil {
		conn.Close()
		return nil, err
	}
	return
}

func (conn DarkConn) Read(data []byte) (n int, err error) {
	buf := make([]byte, len(data), len(data))
	n, err = conn.Conn.Read(buf)
	if n > 0 {
		conn.Decode(buf, data)
	}
	return
}

func (conn DarkConn) Write(data []byte) (n int, err error) {
	buf := make([]byte, len(data), len(data))
	conn.Encode(data, buf)
	n, err = conn.Conn.Write(buf)
	return
}
