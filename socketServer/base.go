package socketServer

//dependencies should saved in vendor directory
//which is under project root, if your IDE ask you to
//download, please enable vendoring
import "github.com/gorilla/websocket"


//create this struct for callback
type connectionElement struct {
	connID uint16
	conn websocket.Conn
}


