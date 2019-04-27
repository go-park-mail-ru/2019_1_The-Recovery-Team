package chat

//easyjson:json
type InitMessagePayload struct {
	SessionID string  `json:"sessionId"`
	Author    *uint64 `json:"author"`
	Receiver  *uint64 `json:"toId"`
	Data      Data    `json:"data"`
}

type SetSessionPayload struct {
	SessionID string `json:"sessionId"`
}

type UpdateMessagePayload struct {
	MessageId *uint64
	Data      Data `json:"data"`
}
