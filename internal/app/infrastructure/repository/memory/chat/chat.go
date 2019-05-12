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

func (c *Chat) broadcast(receiver *uint64, actionType string, payload interface{}) {
	if receiver != nil {
		c.Log.Info("Private message",
			zap.Uint64("user_id", *receiver))
		session, ok := c.Sessions.Load(*receiver)
		if ok {
			c.Log.Info("Private message receiver is online",
				zap.Uint64("user_id", *receiver))
			if user, ok := c.Users.Load(session); ok {
				user.(*chat.User).Messages <- &chat.Action{
					Type:    actionType,
					Payload: payload,
				}
			}
		}
		return
	}

	c.Users.Range(func(key, value interface{}) bool {
		user := value.(*chat.User)
		user.Messages <- &chat.Action{
			Type:    actionType,
			Payload: payload,
		}
		return true
	})
}

func (c *Chat) initGlobalMessage(action *chat.Action) {
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
		return
	}

	user, ok := c.Users.Load(payload.SessionID)
	if !ok {
		c.Log.Error("Can't find user requested global messages")
		return
	}

	user.(*chat.User).Messages <- &chat.Action{
		Type:    chat.SetGlobalMessages,
		Payload: messages,
	}
}

func (c *Chat) initMessage(action *chat.Action) {
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
		return
	}

	c.broadcast(created.Receiver, chat.SetMessage, created)
}

func (c *Chat) initUpdateMessage(action *chat.Action) {
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
		return
	}

	c.broadcast(updated.Receiver, chat.SetUpdateMessage, updated)
}

func (c *Chat) initPrinting(action *chat.Action) {
	payload := action.Payload.(*chat.InitPrintingPayload)
	c.Users.Range(func(key, value interface{}) bool {
		if key.(string) == payload.SessionID {
			return true
		}
		value.(*chat.User).Messages <- &chat.Action{
			Type: chat.SetPrinting,
			Payload: chat.SetPrintingPayload{
				Id: *payload.Author,
			},
		}
		return true
	})
}

func (c *Chat) initDeleteMessage(action *chat.Action) {
	payload := action.Payload.(*chat.InitDeleteMessagePayload)
	message := &chat.Message{
		ID:     payload.Id,
		Author: payload.Author,
	}

	deleted, err := c.MessageManager.Delete(message)
	if err != nil {
		c.Log.Error("Database error",
			zap.String("error", err.Error()))
		return
	}

	c.broadcast(deleted.Receiver, chat.SetDeleteMessage, deleted)
}

func (c *Chat) processAction() {
	for action := range c.Actions {
		c.Log.Info("Receive message")
		switch action.Type {
		case chat.InitGlobalMessages:
			{
				c.initGlobalMessage(action)
			}
		case chat.InitMessage:
			{
				c.initMessage(action)
			}
		case chat.InitUpdateMessage:
			{
				c.initUpdateMessage(action)
			}
		case chat.InitPrinting:
			{
				c.initPrinting(action)
			}
		case chat.InitDeleteMessage:
			{
				c.initDeleteMessage(action)
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
