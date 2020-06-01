package websocket


type UniqueIdentification string



// websocket message
type MessageFormat struct {
	Event       string   `json:"event"`
	From        UniqueIdentification
	Data    	interface{} `json:"data,omitempty"`
}

