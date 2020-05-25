package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

// websocket message
type MessageFormat struct {
	Route       string                 `json:"route"`
	From        UniqueIdentification
	Data    	interface{} `json:"data,omitempty"`
}

type WebSocket struct {
	W http.ResponseWriter
	R *http.Request
	Conn *Connection
	MessageFormat MessageFormat
}

type NoticeController func(*WebSocket)

type socketRoute map[string]NoticeController

func (sr *socketRoute)Add (routeName string,fuc NoticeController) {
	(*sr)[routeName] = fuc
}

var Route = make(socketRoute)

// websocket middleware
func Middleware(w http.ResponseWriter, r *http.Request,auth Auth) {
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
		if err,conn.UniqueIdentification = auth.Identity(w,r);err == nil {
			Manager.Connected(conn.UniqueIdentification,conn)
			_ = auth.ConnDone(conn)
		}else {
			goto ERR
		}
	}
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		for {
			<-ticker.C
			if err = conn.WriteMessage([]byte("ping"));err != nil {
				return
			}
		}
	}()
	defer conn.Close()
	for {
		if data , err = conn.ReadMessage();err != nil {
			goto ERR
		}
		var messageFormat MessageFormat
		if err := json.Unmarshal(data,&messageFormat);err != nil {
			_ = conn.WriteMessage([]byte("data format fail"))
			continue
		}
		messageFormat.From = conn.UniqueIdentification
		if fuc,ok := Route[messageFormat.Route];ok {
			fuc(&WebSocket{
				W:    w,
				R:    r,
				Conn: conn,
				MessageFormat: messageFormat,
			})
		}else {
			_ = conn.WriteMessage([]byte("route not existent"))
		}
	}
	ERR:
		conn.Close()
}


type UniqueIdentification string

type Auth interface {
	Identity(w http.ResponseWriter, r *http.Request) (error,UniqueIdentification)
	ConnDone(conn *Connection) error
}