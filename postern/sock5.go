package postern

import (
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	sockVer = 0x05 // socks5
	noAuth	= 0x00 // success
	cmdType = 0x01 // connect
	rsvByte = 0x00 // meanless
	addType = 0x03 // domain
	succRep = 0x00 // success
)

var (
	errVer = errors.New("version error")
	errCmd = errors.New("command error")
	errRsv = errors.New("reserve error")
	errAdd = errors.New("address error")
	errMtd = errors.New("method error")
)

func HandShake(conn net.Conn) (err error) {
	buf := make([]byte, 257, 257) // 2 + 255
	var n int
	if n, err = io.ReadAtLeast(conn, buf, 2); err != nil {
		return
	}
	// check protocol
	if buf[0] != sockVer {
		return errVer
	}
	if byte(n) != buf[1]+2 {
		return errMtd
	}
	// reply no auth
	_, err = conn.Write([]byte{sockVer, noAuth})
	return nil
}

func DecodeDest(buf []byte) (dest string, err error) {
	if buf[0] != sockVer {
		return "", errVer
	}
	if buf[1] != cmdType {
		return "", errCmd
	}
	if buf[2] != rsvByte {
		return "", errRsv
	}
	// handle domain only, support ip later
	if buf[3] != addType {
		return "", errAdd
	}

	dmLen := int(buf[4])
	domain := string(buf[5 : 5 + dmLen])
	var port int
	port = int(buf[5 + dmLen])
	port = (port << 8) | int(buf[5 + dmLen + 1])
	dest = domain + ":" + strconv.Itoa(port)
	return
}

func LocalReadDest(conn net.Conn) (dest string, err error) {
	msgLen := 1 + 1 + 1 + 1 + 1 + 255 + 2 // ver + cmd + rsv + type + len + domain + port
	buf := make([]byte, msgLen, msgLen)
	io.ReadFull(conn, buf)
	// get dest addr
	dest, err = DecodeDest(buf)
	if err != nil {
		return
	}
	// reply success
	_, err = conn.Write([]byte{sockVer, succRep, rsvByte, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	return "", nil
}

// we encode it ourself, so it has no err
func EncodeDest(dest string) []byte {
	arr := strings.Split(dest, ":")
	domain := arr[0]
	port, _ := strconv.Atoi(arr[1])

	msgLen := 1 + 1 + 1 + len(domain) + 2 // ver + type + len + domain + port
	buf := make([]byte, msgLen, msgLen)

	buf[0] = sockVer
	buf[1] = addType
	buf[2] = byte(len(domain))               // domain length
	copy(buf[2:], []byte(domain))            // domain bytes
	buf[msgLen-2] = byte((port >> 8) & 0xFF) // port high 8 bits
	buf[msgLen-1] = byte(port & 0xFF)        // port low 8 bits

	return buf
}

func RemoteReadDest(conn net.Conn) (dest string, err error) {
	buf := make([]byte, 3, 3)
	io.ReadFull(conn, buf)
	if buf[0] != sockVer {
		return "", errVer
	}
	if buf[1] != addType {
		return "", errAdd
	}
	// read domain & port
	dmLen := int(buf[2])
	buf = make([]byte, dmLen + 2, dmLen + 2)
	io.ReadFull(conn, buf)
	// decode
	domain := string(buf[:dmLen])
	var port int
	port = int(buf[dmLen])
	port = (port << 8) | int(buf[dmLen + 1])
	dest = domain + ":" + strconv.Itoa(port)
	return
}
