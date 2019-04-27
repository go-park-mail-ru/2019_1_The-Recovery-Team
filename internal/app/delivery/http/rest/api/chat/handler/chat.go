package handler

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

func Connect(chatManager *usecase.ChatInteractor, sessionManager *session.SessionClient, log *zap.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Get context data(sessionID, profileID)
		cookie, err := r.Cookie("session_id")
		var profileID *uint64

		if err == nil {
			sessionID := &session.SessionId{
				Id: cookie.Value,
			}

			if response, err := (*sessionManager).Get(context.Background(), sessionID); err == nil {
				profileID = &response.Id
			}
		}

		// Upgrade connection
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		user := &chat.User{
			Id:        profileID,
			SessionID: uuid.NewV4().String(),
			Conn:      conn,
			Log:       log,
			Messages:  make(chan interface{}, 10),
		}

		go chatManager.Run()
		chatManager.Connection() <- user
	}
}
