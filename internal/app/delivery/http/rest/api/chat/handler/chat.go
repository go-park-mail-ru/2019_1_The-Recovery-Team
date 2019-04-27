package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/middleware"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

func Connect(chatManager *usecase.ChatInteractor, log *zap.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Get context data(sessionID, profileID)
		sessionID := r.Context().Value(middleware.SessionID).(string)
		profileID := r.Context().Value(middleware.ProfileID).(uint64)

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
			Id:        &profileID,
			SessionID: sessionID,
			Conn:      conn,
			Log:       log,
		}

		go chatManager.Run()
		chatManager.Connection() <- user
	}
}
