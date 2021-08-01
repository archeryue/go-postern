package postern

import (
	"net"
)

type DarkConn struct {
	net.Conn
	Cipher
}

func NewConn(conn net.Conn, cipher Cipher) *DarkConn {
	return &DarkConn{conn, cipher}
}

func DarkDial(dest, remote string, cipher Cipher) (conn *DarkConn, err error) {
	c, err := net.Dial("tcp", remote)
	if err != nil {
		return
	}
	conn = NewConn(c, cipher)
	if _, err := conn.Write(EncodeDest(dest)); err != nil {
		conn.Close()
		return nil, err
	}
	return
}

// overload net.Conn.Read()
func (conn *DarkConn) Read(buf []byte) (n int, err error) {
	n, err = conn.Conn.Read(buf)
	if n > 0 {
		data := conn.Decrypt(buf[:n])
		copy(buf, data)
	}
	return
}

// overload net.Conn.Write()
func (conn *DarkConn) Write(data []byte) (n int, err error) {
	buf := conn.Encrypt(data)
	n, err = conn.Conn.Write(buf)
	return
}
