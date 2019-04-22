package socketServer

//dependencies should saved in vendor directory
//which is under project root, if your IDE ask you to
//download, please enable vendoring
import (
	"github.com/gobwas/ws"
	"log"
	"net"
)


//create this struct for callback
type connectionElement struct {
	connID uint16
	conn net.Conn
}


func Runner() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	//zero-copy upgrade
	u := ws.Upgrader{
		OnHeader: func(key, value []byte) (err error) {
			log.Printf("non-websocket header: %q=%q", key, value)
			return
		},
	}

	//accept incoming connection
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("connection:", err)
		}
		_, err = u.Upgrade(conn)
		if err != nil {
			// handle error
		}
	}
}
