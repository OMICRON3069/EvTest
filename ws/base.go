package ws

//dependencies should saved in vendor directory
//which is under project root, if your IDE ask you to
//download, please enable vendoring
import (
	"EvTest/evBus"
	"log"
	"net"
)

//Ticket office accepts incoming
//connection and publish it to event bus
func TicketOffice(bus evBus.EventBus) {
	//TODO:accept customize addr
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
		return
	}

	//zero-copy upgrade
	u := Upgrader{
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
			return
		}
		_, err = u.Upgrade(conn)
		if err != nil {
			// handle error
			log.Print("upgrade:",err)
			return
		}
		bus.Publish("socket:incoming",conn)
	}
}
