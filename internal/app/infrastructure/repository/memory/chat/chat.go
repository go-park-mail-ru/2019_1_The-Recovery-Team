package chat

import (
	"sync"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	"go.uber.org/zap"
)

func NewRepo(log *zap.Logger, messageManager *usecase.MessageInteractor) *Chat {
	return &Chat{
		Users:      &sync.Map{},
		Connect:    make(chan *chat.User, 10),
		Disconnect: make(chan *chat.User, 10),
		Actions:    make(chan *chat.Action, 100),

		MessageManager: messageManager,
		Log:            log,
	}
}

type Chat struct {
	Users      *sync.Map
	Connect    chan *chat.User
	Disconnect chan *chat.User
	Actions    chan *chat.Action

	MessageManager *usecase.MessageInteractor
	Log            *zap.Logger
}

func (c *Chat) processAction() {
	for action := range c.Actions {
		switch action.Type {
		case chat.InitMessage:
			{
				payload := action.Payload.(*chat.InitMessagePayload)
				message := &chat.Message{
					Author:   payload.Author,
					Receiver: payload.Receiver,
					Data: chat.Data{
						Text: payload.Data.Text,
					},
				}
				created, err := c.MessageManager.Create(message)
				if err != nil {
					continue
				}

				c.Users.Range(func(key, value interface{}) bool {
					user := value.(*chat.User)
					user.Messages <- &chat.Action{
						Type:    chat.SetMessage,
						Payload: created,
					}
					return true
				})
			}
		}
	}
}

func (c *Chat) Run() {
	go c.processAction()

	for {
		select {
		case user := <-c.Connect:
			{
				user.Actions = c.Actions
				user.Disconnect = c.Disconnect
				c.Users.Store(user.SessionID, user)
				go user.ListenAndSend(c.Log)
			}
		case user := <-c.Disconnect:
			{
				c.Users.Delete(user.SessionID)
			}
		}
	}
}

func (c *Chat) Connection() chan *chat.User {
	return c.Connect
}
