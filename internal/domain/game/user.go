package game

import (
	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
)

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
		_, rawMessage, err := u.Conn.ReadMessage()
		if err == nil {
			raw := &ActionRaw{}
			if err := easyjson.Unmarshal(rawMessage, raw); err != nil {
				continue
			}

			action := &Action{
				Type: raw.Type,
			}

			switch action.Type {
			case InitPlayers:
				{
					payload := &InitPlayersPayload{}
					if err := easyjson.Unmarshal([]byte(raw.Payload), payload); err != nil {
						continue
					}
					action.Payload = payload
				}
			case InitPlayerMove:
				{
					payload := &InitPlayerMovePayload{}
					if err := easyjson.Unmarshal([]byte(raw.Payload), payload); err != nil {
						continue
					}
					action.Payload = payload
				}
			}

			u.Room.Actions <- action
		}
	}
}
