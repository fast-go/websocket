package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

type Subject struct {
	W http.ResponseWriter
	R *http.Request
	Conn *Connection
	Websocket *WebSocket
	MessageFormat MessageFormat
}

func (s *Subject) Send(content []byte) error {
	return s.Websocket.Manager.Send(s.Conn.UniqueIdentification,content)
}

func (s *Subject) SendToUid(uniqueIdentification UniqueIdentification,content []byte) error {
	return s.Websocket.Manager.Send(uniqueIdentification,content)
}

func (s *Subject) Broadcast(content []byte) {
	s.Websocket.Manager.Broadcast(s.Conn,content)
}

func (s *Subject) IsOnline(uniqueIdentification UniqueIdentification) (*Connection,bool){
	return s.Websocket.Manager.IsOnline(uniqueIdentification)
}

type EventFunc func(*Subject)

type socketEventFunc map[string]EventFunc

//register event
func (sr *socketEventFunc)Register(eventName string,fuc EventFunc) {
	(*sr)[eventName] = fuc
}

//detach event
func (sr *socketEventFunc)Detach(eventName string) {
	delete(*sr,eventName)
}

type WebSocket struct {
	Events  socketEventFunc
	Manager ConnManager
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		Events: make(socketEventFunc),
		Manager: ConnManager {
			Online:new(int32),
			connections:new(sync.Map),
		},
	}
}

func (webSocket *WebSocket)Middleware(w http.ResponseWriter, r *http.Request,drive Drive) {
	drive.ConnBefore(w,r)
	var(
		conn *Connection
		data []byte
	)
	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	if conn,err = InitConnection(ws);err != nil {
		goto ERR
	}else {
		if err,conn.UniqueIdentification = drive.Identity(w,r);err == nil {
			webSocket.Manager.Connected(conn.UniqueIdentification,conn)
			drive.ConnDone(conn)
		}else {
			goto ERR
		}
	}
	go drive.Heartbeat(conn)

	defer func() {
		conn.Close()
		webSocket.Manager.DisConnected(conn.UniqueIdentification)
	}()
	for {
		if data , err = conn.ReadMessage();err != nil {
			goto ERR
		}
		var messageFormat MessageFormat
		if err := json.Unmarshal(data,&messageFormat);err != nil {
			_ = conn.WriteMessage([]byte("data format fail"))
			goto ERR
		}
		messageFormat.From = conn.UniqueIdentification

		if fuc,ok := webSocket.Events[messageFormat.Event];ok {
			fuc(&Subject{
				W:    w,
				R:    r,
				Conn: conn,
				MessageFormat: messageFormat,
				Websocket: webSocket,
			})
		}else {
			_ = conn.WriteMessage([]byte("Event not existent"))
		}
	}
	ERR:
		conn.Close()
	    webSocket.Manager.DisConnected(conn.UniqueIdentification)
}


type Drive interface {

	//before connection starts
	ConnBefore(w http.ResponseWriter, r *http.Request)

	//return user unique id
	Identity(w http.ResponseWriter, r *http.Request) (error,UniqueIdentification)

	//heartbeat detection
	Heartbeat(conn *Connection)

	//connection complete
	ConnDone(conn *Connection)

}