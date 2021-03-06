package postern

import (
	"io"
	"log"
	"net"
)

func Forward(in, out net.Conn, end chan byte) {
	buf := make([]byte, 4096)
	for {
		//TODO: timeout
		n, err := in.Read(buf)
		if n > 0 {
			if _, err = out.Write(buf[:n]); err != nil {
				log.Println("forward write error: ", err)
				break
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Println("forward read error: ", err)
			}
			break
		}	
	}
	end <- 1
}
