package main

import (
	"log"
	"net"
	"os"
	"strconv"

	pst "github.com/archeryue/go-postern/postern"
)

func handle(config *pst.Config, conn net.Conn) {
	defer conn.Close() 
	// first handshake: select method
	err := pst.HandShake(conn)
	if err != nil {
		log.Println("handshake error: ", err)
		return
	}
	log.Println("handshake success")
	// second handshake: read true dest from second request
	dest, err := pst.LocalReadDest(conn)
	if err != nil {
		log.Println("read dest error: ", err)
		return
	}
	log.Println("read dest:" + dest)
	// establish DarkConn with remote server
	remoteAddr := config.RemoteIp + ":" + strconv.Itoa(config.RemotePort)
	remote, err := pst.DarkDial(dest, remoteAddr, pst.NewCipher(config.Key))
	if err != nil {
		log.Println("dark dail error: ", err)
		return
	}
	defer remote.Close()
	log.Println("dark dail success")
	// forward data, the DarkConn will encode/decode automatically
	end := make(chan byte, 2)
	go pst.Forward(conn, remote, end)
	go pst.Forward(remote, conn, end)
	<-end
}

func serve(config *pst.Config) {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(config.LocalPort))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("start serving, port : %v\n", config.LocalPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			continue
		}
		go handle(config, conn)
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
