package chat

import (
	"io"

	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type User struct {
	Id         *uint64
	SessionID  string
	Conn       *websocket.Conn
	Log        *zap.Logger
	Messages   chan interface{}
	Disconnect chan *User
	Actions    chan *Action
}

// ListenAndSend start player listening and sending functions
func (u *User) ListenAndSend(log *zap.Logger) {
	u.Log = log.With(
		zap.String("session_id", u.SessionID),
	)
	go u.send()
	go u.listen()
}

func (u *User) send() {
	for message := range u.Messages {
		if err := u.Conn.WriteJSON(message); err != nil {
			u.Log.Info("Stop sending. Error on writing to connection")
			return
		}
	}
	u.Log.Info("Stop sending. Message channel was closed")
}

func (u *User) listen() {
	for {
		raw := &ActionRaw{}

		// Read json from connection
		err := u.Conn.ReadJSON(raw)
		if err != nil {
			u.Log.Info("Incorrect json",
				zap.String("error", err.Error()))
		}

		switch {
		case websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived):
			{
				u.Log.Info("Stop listening. Error on reading from connection")
				u.Disconnect <- u
				return
			}
		case err == io.ErrUnexpectedEOF:
			{
				u.Log.Warn("Received incorrect JSON")
				continue
			}
		case err != nil:
			{
				u.Log.Info("Stop listening. User was disconnected.",
					zap.String("error", err.Error()))
				u.Disconnect <- u
				return
			}
		}

		// Send action to room actions
		action := &Action{
			Type: raw.Type,
		}

		switch action.Type {
		case InitMessage:
			{
				u.Log.Info("Receive message")
				payload := &InitMessagePayload{}
				if err := easyjson.Unmarshal([]byte(raw.Payload), payload); err != nil {
					u.Log.Warn("Invalid message init payload")
					continue
				}

				payload.Author = u.Id
				payload.SessionID = u.SessionID
				action.Payload = payload
			}
		}

		u.Actions <- action
	}
}
