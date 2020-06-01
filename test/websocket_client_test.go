package test

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

	ws, err := websocket.Dial(u.String(), "", "http://"+host+"/" )

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
	for i := 0 ;i < 10000;i++ {
		fmt.Println(strconv.Itoa(i + 1))
		//go func(bb int) {
		//	fmt.Println(client.SendMessage([]byte(`{"route":"test","data":"`+strconv.Itoa(bb + 1)+`"}`)))
		//}(i)
		msg := `{"event":"enter","data":"`+strconv.Itoa(i + 1)+`"}`
		//fmt.Println(msg)
		_ = client.SendMessage([]byte(msg))
	}
	elapsed := time.Since(t1)
	fmt.Println(elapsed)
	time.Sleep(time.Second * 100)
}