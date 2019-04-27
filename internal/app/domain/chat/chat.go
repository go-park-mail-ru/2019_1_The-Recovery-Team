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

//easyjson:json
type SetSessionPayload struct {
	SessionID string `json:"sessionId"`
}

//easyjson:json
type InitUpdateMessagePayload struct {
	Id        uint64  `json:"messageId"`
	SessionID string  `json:"-"`
	Author    *uint64 `json:"authorId"`
	Data      Data    `json:"data"`
}

//easyjson:json
type InitPrintingPayload struct {
	SessionID string `json:"-"`
	Author    uint64 `json:"authorId"`
}

//easyjson:json
type SetPrintingPayload struct {
	Id uint64 `json:"id"`
}

//easyjson:json
type InitDeleteMessagePayload struct {
	Id        uint64  `json:"messageId"`
	SessionID string  `json:"-"`
	Author    *uint64 `json:"authorId"`
}
