package chat

//easyjson:json
type InitMessagePayload struct {
	SessionID string
	Author    *uint64 `json:"author"`
	Receiver  *uint64 `json:"toId"`
	Data      Data    `json:"data"`
}
