package postern

import "net"

func Forward(in, out net.Conn, end chan byte) {
	//TODO: read from in and write into out
	end <- 1
}
