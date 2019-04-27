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
		Sessions:   &sync.Map{},
		Connect:    make(chan *chat.User, 10),
		Disconnect: make(chan *chat.User, 10),
		Actions:    make(chan *chat.Action, 100),

		MessageManager: messageManager,
		Log:            log,
	}
}

type Chat struct {
	Users      *sync.Map
	Sessions   *sync.Map
	Connect    chan *chat.User
	Disconnect chan *chat.User
	Actions    chan *chat.Action

	MessageManager *usecase.MessageInteractor
	Log            *zap.Logger
}

func (c *Chat) processAction() {
	for action := range c.Actions {
		c.Log.Info("Receive message")
		switch action.Type {
		case chat.InitGlobalMessages:
			{
				payload := action.Payload.(*chat.InitGlobalMessagesPayload)
				query := &chat.Query{
					Author:   payload.Author,
					Receiver: payload.Receiver,
					Start:    payload.Start,
					Limit:    payload.Limit,
				}
				if query.Limit == 0 {
					query.Limit = 10
				}
				messages, err := c.MessageManager.GetGlobal(query)
				if err != nil {
					c.Log.Error("Can't get global messages",
						zap.String("error", err.Error()))
					continue
				}

				user, ok := c.Users.Load(payload.SessionID)
				if !ok {
					c.Log.Error("Can't find user requested global messages")
					continue
				}

				user.(*chat.User).Messages <- &chat.Action{
					Type:    chat.SetGlobalMessages,
					Payload: messages,
				}
			}
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
					c.Log.Error("Database error",
						zap.String("error", err.Error()))
					continue
				}

				if created.Receiver != nil {
					c.Log.Info("Private message",
						zap.Uint64("user_id", *created.Receiver))
					session, ok := c.Sessions.Load(*created.Receiver)
					if ok {
						c.Log.Info("Private message receiver is online",
							zap.Uint64("user_id", *created.Receiver))
						if user, ok := c.Users.Load(session); ok {
							user.(*chat.User).Messages <- &chat.Action{
								Type:    chat.SetMessage,
								Payload: created,
							}
						}
					}
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
		case chat.InitUpdateMessage:
			{
				payload := action.Payload.(*chat.InitUpdateMessagePayload)
				message := &chat.Message{
					ID: payload.Id,
					Data: chat.Data{
						Text: payload.Data.Text,
					},
					Author: payload.Author,
				}

				updated, err := c.MessageManager.Update(message)
				if err != nil {
					c.Log.Error("Database error",
						zap.String("error", err.Error()))
					continue
				}

				if updated.Receiver != nil {
					c.Log.Info("Private message update",
						zap.Uint64("user_id", *updated.Receiver))
					session, ok := c.Sessions.Load(*updated.Receiver)
					if ok {
						c.Log.Info("Private message update receiver is online",
							zap.Uint64("user_id", *updated.Receiver))
						if user, ok := c.Users.Load(session); ok {
							user.(*chat.User).Messages <- &chat.Action{
								Type:    chat.SetUpdateMessage,
								Payload: updated,
							}
						}
					}
					continue
				}

				c.Users.Range(func(key, value interface{}) bool {
					user := value.(*chat.User)
					user.Messages <- &chat.Action{
						Type:    chat.SetUpdateMessage,
						Payload: updated,
					}
					return true
				})
			}
		case chat.InitPrinting:
			payload := action.Payload.(*chat.InitPrintingPayload)
			c.Users.Range(func(key, value interface{}) bool {
				if key.(string) == payload.SessionID {
					return true
				}
				value.(*chat.User).Messages <- &chat.Action{
					Type: chat.SetPrinting,
					Payload: chat.SetPrintingPayload{
						Id: payload.Author,
					},
				}
				return true
			})
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
				if user.Id != nil {
					c.Sessions.Store(*user.Id, user.SessionID)
				}
				go user.ListenAndSend(c.Log)
				user.Messages <- &chat.Action{
					Type: chat.SetSession,
					Payload: chat.SetSessionPayload{
						SessionID: user.SessionID,
					},
				}
			}
		case user := <-c.Disconnect:
			{
				close(user.Messages)
				if user.Id != nil {
					c.Sessions.Delete(*user.Id)
				}
				c.Users.Delete(user.SessionID)
			}
		}
	}
}

func (c *Chat) Connection() chan *chat.User {
	return c.Connect
}
