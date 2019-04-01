package game

//easyjson:json
type Message struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
}

type Action struct {
	Type    string
	Payload interface{}
}
