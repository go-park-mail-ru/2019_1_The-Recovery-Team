package handler

import (
	"math/rand"
	"net/http"
	"sadislands/internal/domain/game"
	"sadislands/internal/usecase"

	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/websocket"

	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//func searchGame(conn *websocket.Conn) {
//	t := time.NewTicker(3 * time.Second)
//	for {
//		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
//			fmt.Println("Client disconnected")
//			t.Stop()
//			return
//		}
//
//		iter++
//		<-t.C
//	}
//}

func Search(gameInteractor *usecase.GameInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		player := &game.Player{
			//ID:        rand.Uint64(),
			ID:        rand.Uint64(),
			SessionID: uuid.NewV4().String(),
			Conn:      conn,
			Messages:  make(chan *game.Message, 10),
		}
		go player.Send()

		gameInteractor.Players() <- player
	}
}
