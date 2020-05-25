package test

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
	"websocket"
)

const (
	Single MessageType = iota
	Group
	SysNotify
	OnlineNotify
	OfflineNotify
)

const (
	Text MediaType = iota
	Image
	File
)

type MessageType int

type MediaType int

type Auth struct {

}

func (auth *Auth)Identity(w http.ResponseWriter, r *http.Request) (error,websocket.UniqueIdentification){
	//验证用户身份，返回用户唯一标识
	return nil,"1"
}

func (auth *Auth)ConnDone(c *websocket.Connection) error {
	//todo 可以自行存储链接状态

	return nil
}


var index int64

func TestSocket(t *testing.T)  {
	websocket.Route.Add("test", func(s *websocket.Socket) {
		atomic.AddInt64(&index,1)
		fmt.Println("接收到消息:",s.MessageFormat.Data)
		fmt.Println(atomic.LoadInt64(&index))
		_ = s.Conn.WriteMessage([]byte("服务端正在处理test方法"))
		websocket.Manager.Send("2","发送给用户2的消息")
	})

	//group send message
	websocket.Route.Add("group_send", func(s *websocket.Socket) {
		websocket.Manager.Broadcast(s.Conn,[]byte("This is group send message"))
	})

	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		websocket.Middleware(writer,request,&Auth{})
	})

	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Print(err)
	}
}

