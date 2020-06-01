# websocket

```go
go get github.com/fast-go/websocket
```

### service
```go
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

		_ = s.Send([]byte("Send to yourself"))

		s.IsOnline("1")

		_ = s.SendToUid("1",[]byte("Send to others"))

		s.Broadcast([]byte("broadcast"))
		
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



```


### client

Local test, sending one million requests, about 1m34.139s
```go
import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/url"
	"strconv"
	"testing"
	"time"
)

type Client struct {
	Host string
	Path string
	Conn *websocket.Conn
}

func NewWebsocketClient(host, path string) *Client {
	u := url.URL{Scheme: "ws", Host: host, Path: path}

	ws, err := websocket.Dial(u.String(), "", "http://"+host+"/")

	fmt.Println(err)
	return &Client{
		Host: host,
		Path: path,
		Conn:ws,
	}
}

func (this *Client) SendMessage(body []byte) error {
	_, err := this.Conn.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func TestClient(t *testing.T)  {
	client := NewWebsocketClient("localhost:9090","/")
    
    	t1 := time.Now()
    	for i := 0 ;i < 1000000;i++ {
    		//go func(bb int) {
    		//	fmt.Println(client.SendMessage([]byte(`{"route":"test","data":"`+strconv.Itoa(bb + 1)+`"}`)))
    		//}(i)
    		_ = client.SendMessage([]byte(`{"route":"test","data":"`+strconv.Itoa(i + 1)+`"}`))
    	}
    	elapsed := time.Since(t1)
    	fmt.Println(elapsed)
    	time.Sleep(time.Second * 100)
}
```

#### client send data format
```json
{"event":"enter","data":"hello world"}
```