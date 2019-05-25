package handler

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	gameDomain "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/game"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/middleware"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/metric"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/memory/game"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"go.uber.org/zap"
)

const (
	addr = "127.0.0.1:8080/"
)

func TestSearch(t *testing.T) {
	metric.RegisterTotalRoomsMetric("game_service")
	metric.RegisterTotalPlayersMetric("game_service")

	log, err := zap.NewProduction()
	sessionManager := session.NewClientMock()
	profileManager := profile.NewClientMock()
	gameManager := usecase.NewGameInteractor(game.NewGameRepo(log))
	if err != nil {
		panic(err)
	}
	go gameManager.Run()

	router := httprouter.New()
	router.GET("/", middleware.Authentication(&sessionManager, Search(&profileManager, gameManager)))
	go http.ListenAndServe(":8080", router)
	time.Sleep(3 * time.Second)

	cookies := []*http.Cookie{{Name: "session_id", Value: "AUTHORIZED", Path: "/"}}
	dialer := websocket.DefaultDialer
	u, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	jar, _ := cookiejar.New(nil)
	dialer.Jar = jar
	dialer.Jar.SetCookies(u, cookies)

	connFirst, _, err := dialer.Dial("ws://"+addr, nil)
	if err != nil {
		panic(err)
	}
	defer connFirst.Close()

	time.Sleep(3 * time.Second)

	cookies = []*http.Cookie{{Name: "session_id", Value: "AUTHORIZED_MIRROR", Path: "/"}}
	dialer = websocket.DefaultDialer
	jar, _ = cookiejar.New(nil)
	dialer.Jar = jar
	dialer.Jar.SetCookies(u, cookies)

	connSecond, _, err := dialer.Dial("ws://"+addr, nil)
	if err != nil {
		panic(err)
	}
	defer connSecond.Close()

	payload, _ := gameDomain.InitPlayerReadyPayload{
		PlayerId: 1,
	}.MarshalJSON()
	message := gameDomain.ActionRaw{
		Type:    gameDomain.InitPlayerReady,
		Payload: string(payload),
	}

	err = connFirst.WriteJSON(message)
	assert.Empty(t, err, "Doesn't init player1 ready")

	payload, _ = gameDomain.InitPlayerReadyPayload{
		PlayerId: 2,
	}.MarshalJSON()
	message = gameDomain.ActionRaw{
		Type:    gameDomain.InitPlayerReady,
		Payload: string(payload),
	}

	err = connSecond.WriteJSON(message)
	assert.Empty(t, err, "Doesn't init player2 ready")

	payload, _ = gameDomain.InitPlayerMovePayload{
		PlayerId: 2,
		Move:     "LEFT",
	}.MarshalJSON()
	message = gameDomain.ActionRaw{
		Type:    gameDomain.InitPlayerMove,
		Payload: string(payload),
	}

	err = connSecond.WriteJSON(message)
	assert.Empty(t, err, "Doesn't move player2 ready")

	time.Sleep(time.Second)

	payload, _ = gameDomain.InitItemUsePayload{
		PlayerId: 2,
		ItemType: "LIFEBUOY",
	}.MarshalJSON()
	message = gameDomain.ActionRaw{
		Type:    gameDomain.InitItemUse,
		Payload: string(payload),
	}

	err = connSecond.WriteJSON(message)
	assert.Empty(t, err, "Doesn't use item player2")

	timeout := time.After(6 * time.Second)
	for {
		select {
		case <-timeout:
			{
				return
			}
		default:
			{
				action := &gameDomain.Action{}
				err := connFirst.ReadJSON(action)
				assert.Empty(t, err, "Return error on correct reading")
				fmt.Println(action)
				if action.Type == gameDomain.SetGameOver {
					return
				}
			}
		}
	}
}
