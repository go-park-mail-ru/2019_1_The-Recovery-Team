package chat

type Query struct {
	Start    int     `json:"start"`
	Limit    int     `json:"limit"`
	Author   *uint64 `json:"author"`
	Receiver *uint64 `json:"receiver"`
}

//easyjson:json
type InitMessagePayload struct {
	MessageInfo
	Receiver *uint64 `json:"toId"`
	Data     Data    `json:"data"`
}

//easyjson:json
type InitGlobalMessagesPayload struct {
	MessageInfo
	Start    int     `json:"start"`
	Limit    int     `json:"limit"`
	Receiver *uint64 `json:"toId"`
}

//easyjson:json
type SetSessionPayload struct {
	SessionID string `json:"sessionId"`
}

//easyjson:json
type InitUpdateMessagePayload struct {
	MessageInfo
	Id   uint64 `json:"messageId"`
	Data Data   `json:"data"`
}

//easyjson:json
type InitPrintingPayload struct {
	MessageInfo
}

//easyjson:json
type SetPrintingPayload struct {
	Id uint64 `json:"id"`
}

//easyjson:json
type InitDeleteMessagePayload struct {
	MessageInfo
	Id uint64 `json:"messageId"`
}
