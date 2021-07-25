package postern

import (
	"net"
	"strconv"
	"strings"
)

func HandShake(conn net.Conn) (err error) {
	return nil
}

func ReadRequest(conn net.Conn) (dest string, err error) {
	return "", nil
}

func GenHead(dest string) []byte {
	arr := strings.Split(dest, ":")
	ip := arr[0]
	port, _ := strconv.Atoi(arr[1])

	ipLen := len(ip)
	headLen := 1 + 1 + ipLen + 2 // type + ipLen + ipStr + port
	buf := make([]byte, headLen, headLen)

	buf[0] = 3 // type 3 : domain name
	buf[1] = byte(ipLen)
	copy(buf[2:], ip)
	buf[2+ipLen] = byte((port >> 8) & 0xFF) // high 8 bits
	buf[2+ipLen+1] = byte(port & 0xFF)      // low 8 bits

	return buf
}
