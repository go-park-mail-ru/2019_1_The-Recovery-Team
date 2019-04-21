package handler

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/profile"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/http/rest/middleware"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/domain/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
	"github.com/gorilla/websocket"

	"github.com/julienschmidt/httprouter"
)

func Search(profileManager *profile.ProfileClient, gameManager *usecase.GameInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Get context data(sessionID, profileID)
		sessionID := r.Context().Value(middleware.SessionID).(string)
		profileID := r.Context().Value(middleware.ProfileID).(uint64)

		// Get user information
		prof, err := (*profileManager).Get(context.Background(),
			&profile.GetRequest{
				Id: profileID,
			})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
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

		user := &game.User{
			SessionID: sessionID,
			Conn:      conn,
			Messages:  make(chan interface{}, 10),
			Info: game.Info{
				ID:       prof.Info.Id,
				Nickname: prof.Info.Nickname,
				Rating:   prof.Info.Score.Record,
				Avatar:   prof.Info.Avatar,
			},
			StoppedSending: make(chan interface{}),
		}

		// Add new user to queue for searching game
		gameManager.Players() <- user
	}
}
