package game

import (
	"io"

	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type User struct {
	SessionID      string
	GameSessionID  string
	Conn           *websocket.Conn
	Log            *zap.Logger
	Room           *Room
	Messages       chan interface{}
	Info           Info
	Loser          bool
	StoppedSending chan interface{}
}

//easyjson:json
type Info struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Rating   int64  `json:"rating"`
	Avatar   string `json:"avatar"`
}

// ListenAndSend start player listening and sending functions
func (u *User) ListenAndSend(log *zap.Logger) {
	u.Log = log.With(
		zap.Uint64("user_id", u.Info.ID),
	)
	go u.send()
	go u.listen()
}

func (u *User) send() {
	for {
		select {
		case message := <-u.Messages:
			{
				if err := u.Conn.WriteJSON(message); err != nil {
					u.Log.Info("Stop sending. Error on writing to connection")
					return
				}
			}
		case <-u.Room.Ctx.Done():
			{
				u.Log.Info("Correct stopping sending")
				close(u.StoppedSending)
				return
			}
		}
	}
}

func (u *User) listen() {
	for {
		raw := &ActionRaw{}

		// Read json from connection
		err := u.Conn.ReadJSON(raw)
		select {
		case <-u.Room.Ctx.Done():
			{
				u.Log.Info("Correct stopping listening")
				return
			}
		default:
		}

		switch {
		case websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived):
			{
				u.Log.Info("Stop listening. Error on reading from connection")

				if u.Room == nil {
					return
				}

				if !u.Room.Closing.Load() {
					u.Room.Closing.Store(true)
					go u.Room.Close(&Action{
						Type:    SetUserDisconnected,
						Payload: u,
					})
				}

				if u.Room.EngineStarted.Load() {
					u.Room.Actions <- &Action{
						Type: InitEngineStop,
					}
				}
				return
			}
		case err == io.ErrUnexpectedEOF:
			{
				u.Log.Warn("Received incorrect JSON")
				continue
			}
		case err != nil:
			{
				u.Log.Info("Stop listening. User was disconnected")
				return
			}
		}

		// Send action to room actions
		action := &Action{
			Type: raw.Type,
		}

		switch action.Type {
		case InitPing:
			{
				u.Messages <- &Action{
					Type: SetPong,
				}
				continue
			}
		case InitPlayers:
			{
				payload := &InitPlayersPayload{}
				if err := easyjson.Unmarshal([]byte(raw.Payload), payload); err != nil {
					u.Log.Warn("Invalid players init payload")
					continue
				}
				action.Payload = payload
			}
		case InitPlayerReady:
			{
				payload := &InitPlayerReadyPayload{}
				if err := easyjson.Unmarshal([]byte(raw.Payload), payload); err != nil {
					u.Log.Warn("Invalid player ready payload")
					continue
				}

				if payload.PlayerId != u.Info.ID {
					continue
				}

				action.Payload = payload
			}
		case InitPlayerMove:
			{
				payload := &InitPlayerMovePayload{}
				if err := easyjson.Unmarshal([]byte(raw.Payload), payload); err != nil {
					u.Log.Warn("Invalid player move payload")
					continue
				}

				if payload.PlayerId != u.Info.ID {
					u.Log.Warn("Player trying to move opponent")
					continue
				}

				action.Payload = payload
			}
		case InitItemUse:
			{
				payload := &InitItemUsePayload{}
				if err := easyjson.Unmarshal([]byte(raw.Payload), payload); err != nil {
					u.Log.Warn("Invalid item use payload")
					continue
				}

				if payload.PlayerId != u.Info.ID {
					u.Log.Warn("Player trying to use opponent item")
					continue
				}

				action.Payload = payload
			}
		default:
			continue
		}

		if !u.Room.EngineStarted.Load() {
			continue
		}

		u.Room.Actions <- action
	}
}
