package chat

const (
	SetMessage  = "SET_CHAT_MESSAGE"
	InitMessage = "INIT_CHAT_MESSAGE"
	InitPing    = "INIT_PING"
	SetPong     = "SET_PONG"
)

//easyjson:json
type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

//easyjson:json
type ActionRaw struct {
	Type    string `json:"type"`
	Payload string `json:"payload,omitempty"`
}
