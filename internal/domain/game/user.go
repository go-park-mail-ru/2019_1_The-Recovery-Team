package game

import "github.com/gorilla/websocket"

type User struct {
	SessionID     string
	GameSessionID string
	Conn          *websocket.Conn
	Room          *Room
	Messages      chan interface{}
	Info          Info
}

//easyjson:json
type Info struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Rating   int    `json:"rating"`
	Avatar   string `json:"avatar"`
}

func (u *User) Send() {
	for {
		select {
		case message := <-u.Messages:
			{
				// TODO: Process error
				if err := u.Conn.WriteJSON(message); err != nil {
					return
				}
			}
		}
	}
}

func (u *User) Listen() {
	for {
		_, raw, err := u.Conn.ReadMessage()
		if err == nil {
			u.Room.Actions <- string(raw)
		}
	}
}
