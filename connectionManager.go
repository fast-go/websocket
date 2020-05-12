package socket

import (
	"fmt"
	"sync"
	"sync/atomic"
)


type ConnManager struct {
	Online *int32
	connections *sync.Map
}

var Manager = ConnManager {
	Online:new(int32),
	connections:new(sync.Map),
}

func (m *ConnManager) Connected(k,v interface{}) {
	m.connections.Store(k,v)
	atomic.AddInt32(m.Online,1)
}

// remove websocket connection by key
// online number - 1
func (m *ConnManager) DisConnected(k interface{}) {
	m.connections.Delete(k)
	atomic.AddInt32(m.Online, -1)
}

// get websocket connection by key
func (m *ConnManager) Get(k interface{}) (v interface{}, ok bool) {
	return m.connections.Load(k)
}

// iter websocket connections
func (m *ConnManager) Foreach(f func(k, v interface{})) {
	m.connections.Range(func(k, v interface{}) bool {
		f(k, v)
		return true
	})
}

// send message to one websocket connection
func (m *ConnManager) Send(k string, msg string) {
	if v, ok := m.Get(k); ok {
		if conn, ok := v.(*Connection); ok {
			if err := conn.WriteMessage([]byte(msg)); err != nil {
				fmt.Println("Send msg error: ", err)
			}
		} else {
			fmt.Println("invalid type, expect *websocket.Conn")
		}
	} else {
		fmt.Println("connection not exist")
	}
}

// send message to multi websocket connections
func (m *ConnManager) SendMulti(keys []string, msg string) {
	for _, k := range keys {
		v, ok := m.Get(k)
		if ok {
			if conn, ok := v.(*Connection); ok {
				if err := conn.WriteMessage([]byte(msg)); err != nil {
					fmt.Println("Send msg error: ", err)
				}
			} else {
				fmt.Println("invalid type, expect *websocket.Conn")
			}
		} else {
			fmt.Println("connection not exist")
		}
	}
}

// broadcast message to all websocket connections otherwise own connection
func (m *ConnManager) Broadcast(conn *Connection, msg string) {
	m.Foreach(func(k, v interface{}) {
		fmt.Println(k)
		if c, ok := v.(*Connection); ok && c != conn {
			fmt.Println("===========")
			if err := c.WriteMessage([]byte(msg)); err != nil {
				fmt.Println("Send msg error: ", err)
			}
		}
	})
}

