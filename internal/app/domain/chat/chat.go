package chat

type Query struct {
	Start    int     `json:"start"`
	Limit    int     `json:"limit"`
	Author   *uint64 `json:"author"`
	Receiver *uint64 `json:"receiver"`
}

//easyjson:json
type InitMessagePayload struct {
	SessionID string  `json:"sessionId"`
	Author    *uint64 `json:"author"`
	Receiver  *uint64 `json:"toId"`
	Data      Data    `json:"data"`
}

//easyjson:json
type InitGlobalMessagesPayload struct {
	Start     int     `json:"start"`
	Limit     int     `json:"limit"`
	SessionID string  `json:"-"`
	Author    *uint64 `json:"authorId"`
	Receiver  *uint64 `json:"toId"`
}

type SetSessionPayload struct {
	SessionID string `json:"sessionId"`
}

type UpdateMessagePayload struct {
	MessageId *uint64
	Data      Data `json:"data"`
}
