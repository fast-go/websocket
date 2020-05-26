# websocket


```go
go get github.com/fast-go/websocket
```

### service
```go

    type Auth struct {}
    
    //Authentication, return user unique ID
    func (auth *Auth)Identity(w http.ResponseWriter, r *http.Request) (error,socket.UniqueIdentification){
        return nil,"1"
    }
    
    //Connection status can be managed by yourself
    //You can manage the links in groups. When sending group messages
    //you only need to traverse the links of the specified group
    func (auth *Auth)ConnDone(c *socket.Connection) error {
        return nil
    }

    //Define route
    socket.Route.Add("/test", func(s *socket.Socket) {
		fmt.Println("runing test action")
		_ = s.Conn.WriteMessage([]byte("send message"))
		fmt.Println(s.MessageFormat)
        websocket.Manager.Send("2","Send message to user 2")
	})
    
    //Middleware connect websocket
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
         socket.Middleware(writer,request,&Auth{})
	})

	// start http service
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Print(err)
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
{"route":"/test","data":"hello world"}
```