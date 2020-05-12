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
    func (auth *Auth)ConnDone(c *socket.Connection) error {
        return nil
    }
    //Define route
    socket.Route.Add("/test", func(s *socket.Socket) {
		fmt.Println("runing test action")
		_ = s.Conn.WriteMessage([]byte("send message"))
		fmt.Println(s.MessageFormat)
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		 //Middleware connect websocket
         socket.Middleware(writer,request,&Auth{})
	})

	// start http service
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Print(err)
	}
```


### client

#### client send data format
```js
{"route":"/test","data":"hello world"}
```