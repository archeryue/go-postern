package postern

import (
	"net"
	"strconv"
	"strings"
)

type DarkConn struct {
	net.Conn
	*Cipher
}

func NewConn(conn net.Conn, cipher *Cipher) *DarkConn {
	return &DarkConn{conn, cipher}
}

func sock5head(dest string) []byte {
	arr := strings.Split(dest, ":")
	ip := arr[0]
	port, _ := strconv.Atoi(arr[1])

	ipLen := len(ip)
	headLen := 1 + 1 + ipLen + 2 // type + ipLen + ipStr + port
	buf := make([]byte, headLen, headLen)

	buf[0] = 3 // type 3 : domain name
	buf[1] = byte(ipLen)
	copy(buf[2:], ip)
	buf[2 + ipLen] = byte((port >> 8) & 0xFF) // high 8 bits
	buf[2 + ipLen + 1] = byte(port & 0xFF) // low 8 bits

	return buf
}

func DarkDial(dest, remote string, cipher *Cipher) (conn *DarkConn, err error) {
	c, err := net.Dial("tcp", remote)
	if (err != nil) {
		return
	}
	conn = NewConn(c, cipher)
	if _, err := conn.Write(sock5head(dest)); err != nil {
		conn.Close()
		return nil, err
	}
	return
}

func (conn DarkConn) Read(data []byte) (n int, err error) {
	buf := make([]byte, len(data), len(data))
	n, err = conn.Conn.Read(buf)
	if (n > 0) {
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
