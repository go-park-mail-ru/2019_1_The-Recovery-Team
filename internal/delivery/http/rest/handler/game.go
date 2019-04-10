package handler

import (
	"net/http"
	"sadislands/internal/delivery/http/rest/middleware"
	"sadislands/internal/domain/game"
	"sadislands/internal/usecase"

	"github.com/gorilla/websocket"

	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Search(profileInteractor *usecase.ProfileInteractor, gameInteractor *usecase.GameInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		sessionID := r.Context().Value(middleware.SessionID).(string)
		profileID := r.Context().Value(middleware.ProfileID).(uint64)

		profile, err := profileInteractor.GetProfile(profileID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		user := &game.User{
			SessionID:    sessionID,
			Conn:         conn,
			Messages:     make(chan interface{}, 10),
			Info: game.Info{
				ID:       profile.ID,
				Nickname: profile.Nickname,
				Rating:   profile.Record,
				Avatar:   profile.Avatar,
			},
		}

		gameInteractor.Players() <- user
	}
}
