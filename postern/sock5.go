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
	noAuth  = 0x00 // success
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
	dmLen := len(buf) - 2
	domain := string(buf[:dmLen])
	var port int
	port = int(buf[dmLen])
	port = (port << 8) | int(buf[dmLen+1])
	dest = domain + ":" + strconv.Itoa(port)
	return
}

func LocalReadDest(conn net.Conn) (dest string, err error) {
	headLen := 1 + 1 + 1 + 1 + 1 // ver + cmd + rsv + type + len
	buf := make([]byte, headLen, headLen)
	io.ReadFull(conn, buf)
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
	buf = make([]byte, dmLen+2, dmLen+2)
	io.ReadFull(conn, buf)
	// get dest addr
	dest, err = DecodeDest(buf)
	if err != nil {
		return
	}
	// reply success
	_, err = conn.Write([]byte{sockVer, succRep, rsvByte, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	return
}

func EncodeDest(dest string) []byte {
	arr := strings.Split(dest, ":")
	domain := arr[0]
	port, _ := strconv.Atoi(arr[1])

	dtLen := len(domain) + 2
	buf := make([]byte, dtLen, dtLen)

	copy(buf, []byte(domain))
	buf[dtLen-2] = byte((port >> 8) & 0xFF) // port high 8 bits
	buf[dtLen-1] = byte(port & 0xFF)        // port low 8 bits
	return buf
}

// we encode it ourself, so it has no err
func EncodeRequest(dest string) []byte {
	dt := EncodeDest(dest)
	msgLen := 1 + 1 + 1 + len(dt) // ver + type + len + domain + port
	buf := make([]byte, msgLen, msgLen)

	buf[0] = sockVer
	buf[1] = addType
	buf[2] = byte(len(dt) - 2) // domain length
	copy(buf[3:], dt)          // domain & port bytes
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
	buf = make([]byte, dmLen+2, dmLen+2)
	io.ReadFull(conn, buf)
	// decode
	dest, err = DecodeDest(buf)
	return
}
