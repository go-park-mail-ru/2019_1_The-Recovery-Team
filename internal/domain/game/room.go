package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

type Room struct {
	ID      string
	Players *sync.Map
	Total   int
	Exclude chan *Player
	Ctx     context.Context
	Cancel  context.CancelFunc

	Log *zap.Logger
}

func NewRoom(log *zap.Logger, closed chan *Room) *Room {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "closed", closed)
	id := uuid.NewV4().String()
	return &Room{
		ID:      id,
		Players: &sync.Map{},
		Exclude: make(chan *Player, 1),
		Ctx:     ctx,
		Cancel:  cancel,
		Log: log.With(
			zap.String("room_id", id),
		),
	}
}

func (r *Room) Run() {
	players := make([]*Player, 0, 2)

	r.Players.Range(func(key, value interface{}) bool {
		player := value.(*Player)
		players = append(players, player)
		go player.Listen()
		return true
	})

	players[0].Messages <- &Message{
		Status:  200,
		Payload: fmt.Sprintf("Your opponent %d", players[1].ID),
	}

	players[1].Messages <- &Message{
		Status:  200,
		Payload: fmt.Sprintf("Your opponent %d", players[0].ID),
	}

	for {
		select {
		case player := <-r.Exclude:
			{
				action := &Action{
					Type:    "Disconnect",
					Payload: player,
				}
				r.Close(action)
				return
			}
		}
	}
}

func (r *Room) Broadcast(message *Message) error {
	var err error

	r.Players.Range(func(key, value interface{}) bool {
		mes, _ := message.MarshalJSON()
		err = value.(*Player).Conn.WriteMessage(websocket.TextMessage, mes)
		if err != nil {
			return false
		}
		return true
	})

	return err
}

func (r *Room) Close(action *Action) {
	switch action.Type {
	case "Disconnect":
		{
			leaver := action.Payload.(*Player)
			var player *Player

			r.Players.Range(func(key, value interface{}) bool {
				player = value.(*Player)
				if player.ID != leaver.ID {
					return false
				}
				return true
			})

			player.Messages <- &Message{
				Status:  200,
				Payload: "Your opponent left the game(",
			}

			r.Log.Info("User left the game",
				zap.Uint64("leaver_id", leaver.ID),
				zap.Uint64("opponent_id", player.ID),
				zap.String("room_id", r.ID))
		}
	}

	r.Cancel()
	r.Players.Range(func(key, value interface{}) bool {
		player := value.(*Player)

		player.Conn.SetReadDeadline(time.Now().Add(time.Second))
		player.Conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
		player.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		player.Conn.Close()
		r.Log.Info("User disconnected from game",
			zap.Uint64("user_id", player.ID))
		return true
	})

	r.Log.Info("Room closed")
}
