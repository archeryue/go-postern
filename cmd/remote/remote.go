package main

import (
	"log"
	"net"
	"os"
	"strconv"

	pst "github.com/archeryue/go-postern/postern"
)

func handle(conn *pst.DarkConn) {
	defer conn.Close()
	// read dest addr
	dest, err := pst.RemoteReadDest(conn)
	if err != nil {
		log.Println("read dest error: ", err)
		return
	}
	// connect dest host
	destConn, err := net.Dial("tcp", dest)
	if err != nil {
		log.Println("dail target error: ", err)
		return
	}
	defer destConn.Close()
	// forward data, the DarkConn(conn) will encode/decode automatically
	end := make(chan byte, 2)
	go pst.Forward(conn, destConn, end)
	go pst.Forward(destConn, conn, end)
	// each one of two conns is end
	<-end
}

func serve(config *pst.Config) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(config.RemotePort))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("start serving, port : %v\n", config.RemotePort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			continue
		}
		cipher := pst.NewCipher(config.Key, config.Method)
		go handle(pst.NewConn(conn, cipher))
	}
}

func main() {
	path := os.Args[1]
	config, err := pst.LoadConfig(path)
	if err != nil {
		log.Printf("Config Error: %s\n", path)
		return
	}
	serve(config)
}
