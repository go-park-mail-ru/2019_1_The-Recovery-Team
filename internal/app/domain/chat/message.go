package chat

import "time"

//easyjson:json
type Message struct {
	ID       uint64    `json:"messageId" example:"1"`
	Author   *uint64   `json:"authorId" example:"2"`
	Receiver *uint64   `json:"toId" example:"3"`
	Created  time.Time `json:"created"`
	Edited   bool      `json:"isEdited"`
	Data     Data      `json:"data"`
}

//easyjson:json
type Data struct {
	Text string `json:"text"`
}
