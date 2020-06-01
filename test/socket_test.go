package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"
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

func (auth *Auth) ConnBefore(w http.ResponseWriter, r *http.Request) {

}

func (auth *Auth)ConnDone(c *websocket.Connection) {
	//todo 可以自行存储链接状态,可对链接进行分组管理.组消息发送的时候只需要遍历制定组的链接
}

func (auth *Auth)Heartbeat(c *websocket.Connection) {
	ticker := time.NewTicker(time.Second * 2)
	for {
		<-ticker.C
		if err := c.WriteMessage([]byte("ping"));err != nil {
			return
		}
	}
}

func TestSocket(t *testing.T)  {
	im := websocket.NewWebSocket()
	im.Events.Register("enter", func(s *websocket.Subject) {
		fmt.Println("message:",s.MessageFormat.Data)
		_ = s.Send([]byte("Send to yourself"))

		//s.IsOnline("1")

		//_ = s.SendToUid("1",[]byte("Send to others"))

		//s.Broadcast([]byte("broadcast"))

	})

	//group send message
	im.Events.Register("group_send", func(s *websocket.Subject) {
		im.Manager.Broadcast(s.Conn,[]byte("This is group send message"))
	})

	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		im.Middleware(writer,request,&Auth{})
	})


	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Print(err)
	}
}

