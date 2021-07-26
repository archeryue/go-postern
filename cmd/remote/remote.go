package remote

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
	// connect target
	target, err := net.Dial("tcp", dest)
	if err != nil {
		log.Println("dail target error: ", err)
		return
	}
	defer target.Close()
	// forward data, the DarkConn will encode/decode automatically
	end := make(chan byte, 2)
	go pst.Forward(conn, target, end)
	go pst.Forward(target, conn, end)
	<-end
}

func serve(config *pst.Config) {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(config.RemotePort))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("start serving, port : %v\n", config.RemotePort)

	cipher := pst.NewCipher(config.Key)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			continue
		}
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
