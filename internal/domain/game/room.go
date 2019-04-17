package game

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

//easyjson:json
type Message struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
}

type Room struct {
	ID            string
	Users         *sync.Map
	Total         atomic.Uint64
	EngineStarted atomic.Bool
	Closed        chan *Room
	Actions       chan *Action
	EngineStopped chan interface{}

	Ctx     context.Context
	Cancel  context.CancelFunc
	Closing atomic.Bool

	Log *zap.Logger
}

// NewRoom creates new instance of room
func NewRoom(log *zap.Logger, closed chan *Room) *Room {
	id := uuid.NewV4().String()
	ctx, cancel := context.WithCancel(context.Background())

	room := &Room{
		ID:      id,
		Users:   &sync.Map{},
		Closed:  closed,
		Actions: make(chan *Action, 10),
		Log: log.With(
			zap.String("room_id", id),
		),
		EngineStopped: make(chan interface{}),

		Ctx:     ctx,
		Cancel:  cancel,
		Closing: atomic.Bool{},
	}
	room.Closing.Store(false)
	room.EngineStarted.Store(false)
	return room
}

// Run starts game engine at room
func (r *Room) Run(sendInto func(action interface{})) {
	users := make([]*User, 0, 2)

	// Get list of users
	r.Users.Range(func(key, value interface{}) bool {
		user := value.(*User)
		users = append(users, user)
		return true
	})

	// Send information about opponents
	info := users[1].Info
	users[0].Messages <- &Action{
		Type:    SetOpponent,
		Payload: info,
	}

	info = users[0].Info
	users[1].Messages <- &Action{
		Type:    SetOpponent,
		Payload: info,
	}

	// Init players id in engine
	playersId := make([]uint64, 0, len(users))
	for _, user := range users {
		playersId = append(playersId, user.Info.ID)
	}

	payload := &InitPlayersPayload{
		PlayersId: playersId,
	}

	action := &Action{
		Type:    InitPlayers,
		Payload: payload,
	}
	sendInto(action)

	r.EngineStarted.Store(true)

	// Receive actions form users
	for {
		select {
		case action := <-r.Actions:
			{
				sendInto(action)
				if action.Type == InitEngineStop {
					r.Log.Info("Action channel was closed")
					return
				}
			}
		case <-r.Ctx.Done():
			{
				r.Log.Info("Action channel was closed by running flag")
				return
			}
		}
	}
}

// Broadcast sends messages to all users in room
func (r *Room) Broadcast(action *Action) {
	r.Users.Range(func(key, value interface{}) bool {
		value.(*User).Messages <- action
		return true
	})
}

// Close removes room and stops engine
func (r *Room) Close(action *Action) {
	if r.Closing.Load() {
		return
	}
	r.Closing.Store(true)
	r.Log.Info("Closing room")

	// Stop running engine
	if r.EngineStarted.Load() {
		r.Log.Info("Stopping engine")
		<-r.EngineStopped
	}
	r.Cancel()

	switch action.Type {
	case SetUserDisconnected:
		{
			// Send information about leaver to players
			leaver := action.Payload.(*User)
			var opponent User

			r.Users.Range(func(key, value interface{}) bool {
				user := value.(*User)
				close(user.Messages)
				if user.Info.ID != leaver.Info.ID {
					opponent = *user
					<-user.StoppedSending
					user.Conn.WriteJSON(&Action{
						Type: SetOpponentLeave,
					})
					user.Conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
					user.Conn.WriteMessage(websocket.CloseMessage,
						websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					time.Sleep(1 * time.Second)
					user.Conn.Close()

					r.Log.Info("User disconnected from game",
						zap.Uint64("user_id", user.Info.ID))
				}

				return true
			})

			r.Log.Info("User left the game",
				zap.Uint64("leaver_id", leaver.Info.ID),
				zap.Uint64("opponent_id", opponent.Info.ID),
				zap.String("room_id", r.ID))
		}
	case SetGameOver:
		{
			// Close player connections
			r.Users.Range(func(key, value interface{}) bool {
				user := value.(*User)
				close(user.Messages)

				<-user.StoppedSending
				user.Conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				user.Conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				time.Sleep(1 * time.Second)
				user.Conn.Close()

				r.Log.Info("User disconnected from game",
					zap.Uint64("user_id", user.Info.ID))
				return true
			})
		}
	}

	close(r.Actions)

	r.Closed <- r
}

// ActionCallback process actions received from engine
func (r *Room) ActionCallback(action *Action) {
	switch action.Type {
	case SetEngineStop:
		{
			r.Log.Info("Engine was stopped")
			select {
			case <-r.EngineStopped:
				{
					return
				}
			default:
				{
					r.Log.Info("Setting flag for engine stopped")
					close(r.EngineStopped)
					return
				}
			}
		}
	case SetGameOver:
		{
			r.Broadcast(action)
			go r.Close(&Action{
				Type: SetGameOver,
			})
			return
		}
	}
	r.Broadcast(action)
}
