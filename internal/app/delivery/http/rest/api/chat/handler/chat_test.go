package handler

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	messageRepo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql/message"

	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	chatRepo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/memory/chat"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"go.uber.org/zap"
)

const (
	addr = "127.0.0.1:8080/"
)

func TestConnect(t *testing.T) {
	log, err := zap.NewProduction()
	messageManager := usecase.NewMessageInteractor(&messageRepo.RepoMock{})
	chatManager := usecase.NewChatInteractor(chatRepo.NewRepo(log, messageManager))
	sessionManager := session.NewClientMock()
	if err != nil {
		panic(err)
	}
	go chatManager.Run()

	router := httprouter.New()
	router.GET("/", Connect(chatManager, &sessionManager, log))
	go http.ListenAndServe(":8080", router)
	time.Sleep(1 * time.Second)

	cookies := []*http.Cookie{{Name: "session_id", Value: "AUTHORIZED", Path: "/"}}
	dialer := websocket.DefaultDialer
	u, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	jar, _ := cookiejar.New(nil)
	dialer.Jar = jar
	dialer.Jar.SetCookies(u, cookies)

	ws, _, err := dialer.Dial("ws://"+addr, nil)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	action := &chat.Action{}

	// Receive session id
	err = ws.ReadJSON(action)
	assert.Empty(t, err, "Doesn't read correct chat session")

	message := chat.ActionRaw{
		Type: chat.InitPing,
	}

	// Send ping
	err = ws.WriteJSON(message)
	assert.Empty(t, err, "Doesn't send ping")

	// Receive pong
	err = ws.ReadJSON(action)
	assert.Empty(t, err, "Doesn't receive pong")

	payload, _ := chat.InitMessagePayload{
		Data: chat.Data{
			Text: "text",
		},
	}.MarshalJSON()
	message = chat.ActionRaw{
		Type:    chat.InitMessage,
		Payload: string(payload),
	}

	// Send global message
	err = ws.WriteJSON(message)
	assert.Empty(t, err, "Doesn't send correct message")

	// Receive global message
	err = ws.ReadJSON(action)
	assert.Empty(t, err, "Doesn't receive correct message")

	messageId := uint64(action.Payload.(map[string]interface{})["messageId"].(float64))
	payload, _ = chat.InitUpdateMessagePayload{
		Id: messageId,
		Data: chat.Data{
			Text: "updated",
		},
	}.MarshalJSON()
	message = chat.ActionRaw{
		Type:    chat.InitUpdateMessage,
		Payload: string(payload),
	}

	// Update message
	err = ws.WriteJSON(message)
	assert.Empty(t, err, "Doesn't update correct message")

	// Receive updated message
	err = ws.ReadJSON(action)
	assert.Empty(t, err, "Doesn't receive updated global message")

	payload, _ = chat.InitDeleteMessagePayload{
		Id: messageId,
	}.MarshalJSON()
	message = chat.ActionRaw{
		Type:    chat.InitDeleteMessage,
		Payload: string(payload),
	}

	// Delete message
	err = ws.WriteJSON(message)
	assert.Empty(t, err, "Doesn't delete correct message")

	// Receive deleted message
	err = ws.ReadJSON(action)
	assert.Empty(t, err, "Doesn't receive deleted global message")

	message = chat.ActionRaw{
		Type: chat.InitPrinting,
	}

	// Start printing
	err = ws.WriteJSON(message)
	assert.Empty(t, err, "Doesn't delete correct message")

	payload, _ = chat.InitGlobalMessagesPayload{}.MarshalJSON()
	message = chat.ActionRaw{
		Type:    chat.InitGlobalMessages,
		Payload: string(payload),
	}

	// Request global messages
	err = ws.WriteJSON(message)
	assert.Empty(t, err, "Doesn't request global messages")

	// Receive global messages
	err = ws.ReadJSON(action)
	assert.Empty(t, err, "Doesn't receive global messages")
}
