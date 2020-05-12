package test

import (
	"fmt"
	"net/http"
	"socket"
	"testing"
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

func (auth *Auth)Identity(w http.ResponseWriter, r *http.Request) (error,socket.UniqueIdentification){
	//验证用户身份
	return nil,"1"
}

func (auth *Auth)ConnDone(c *socket.Connection) error {
	//todo 可以自行存储链接状态

	return nil
}

func TestSocket(t *testing.T)  {
	socket.Route.Add("/test", func(s *socket.Socket) {
		fmt.Println("正在处理test方法")
		_ = s.Conn.WriteMessage([]byte("服务端正在处理test方法"))
		fmt.Println(s.MessageFormat)
	})

	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		socket.Middleware(writer,request,&Auth{})
	})

	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Print(err)
	}
}


