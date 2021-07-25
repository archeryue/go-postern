package local

import (
	"log"
	"net"
	"os"
	"strconv"

	pst "github.com/archeryue/go-postern/postern"
)

func handle(config *pst.Config, conn net.Conn) {
	defer conn.Close() 
	//TODO: sock5 handshake, get local request
	//TODO: send request to remote, get response
	//TODO: send response back to local
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
