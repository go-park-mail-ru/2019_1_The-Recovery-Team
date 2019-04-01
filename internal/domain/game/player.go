package game

import "github.com/gorilla/websocket"

type Player struct {
	ID            uint64
	SessionID     string
	GameSessionID string
	Conn          *websocket.Conn
	Room          *Room
	Messages      chan *Message
}

func (p *Player) Send() {
	for {
		select {
		case message := <-p.Messages:
			{
				msg, _ := message.MarshalJSON()

				// TODO: Process error
				if err := p.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					return
				}
			}
		}
	}
}

func (p *Player) Listen() {
	for {
		_, raw, err := p.Conn.ReadMessage()
		if err == nil {
			message := string(raw)
			switch message {
			case "exit":
				{
					p.Room.Exclude <- p
					return
				}

			}
		}
		return
	}
}
