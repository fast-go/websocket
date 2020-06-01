package websocket

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)


type ConnManager struct {
	Online *int32
	connections *sync.Map
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

// check if the connection exists
func (m *ConnManager)IsOnline(uniqueIdentification UniqueIdentification) (*Connection,bool) {
	v,b := m.Get(uniqueIdentification)
	if b {
		if c ,ok := v.(*Connection);ok{
			return c,ok
		}
	}
	return nil,false
}

// send message to one websocket connection
func (m *ConnManager) Send(k UniqueIdentification, msg []byte) (err error) {
	if v, ok := m.Get(k); ok {
		if conn, ok := v.(*Connection); ok {
			if err = conn.WriteMessage(msg); err != nil {
				return err
			}
		} else {
			err = errors.New("invalid type, expect *websocket.Conn")
		}
	} else {
		err = errors.New("connection not exist")
	}
	return err
}

// send message to multi websocket connections
func (m *ConnManager) SendMulti(keys []UniqueIdentification, msg string) {
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
func (m *ConnManager) Broadcast(conn *Connection, msg []byte) {
	m.Foreach(func(k, v interface{}) {
		if c, ok := v.(*Connection); ok && c != conn {
			if err := c.WriteMessage(msg); err != nil {
				fmt.Println("Send msg error: ", err)
			}
		}
	})
}

