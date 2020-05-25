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
	for i := 0 ;i < 200000;i++ {
		go func(bb int) {
			fmt.Println(client.SendMessage([]byte(`{"route":"test","data":"`+strconv.Itoa(bb)+`"}`)))
		}(i)
	}
	time.Sleep(time.Second * 100)
}