package chat

import (
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
	uuid "github.com/satori/go.uuid"

	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql/message"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"go.uber.org/zap"
)

var (
	receiverId        uint64 = 2
	receiverSessionId        = uuid.NewV4().String()
)

func repo() *Chat {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	messageManager := usecase.NewMessageInteractor(&message.RepoMock{})
	return NewRepo(log, messageManager)
}

func TestNewRepo(t *testing.T) {
	assert.NotEmpty(t, repo(), "Empty chat on correct data")
}

func TestBroadcast(t *testing.T) {
	repo := repo()

	for i := 1; i < 3; i++ {
		id := uint64(i)
		sessionID := uuid.NewV4().String()
		repo.Users.Store(sessionID, &chat.User{
			Id:        &id,
			SessionID: sessionID,
			Messages:  make(chan interface{}, 5),
		})

		repo.Sessions.Store(id, sessionID)
	}

	//Test private message
	receiver := uint64(2)
	repo.broadcast(&receiver, chat.InitMessage, nil)
	sessionId, _ := repo.Sessions.Load(receiver)
	user, _ := repo.Users.Load(sessionId)
	close(user.(*chat.User).Messages)
	_, ok := <-user.(*chat.User).Messages
	assert.Equal(t, true, ok, "Doesn't receive private messsage")

	user.(*chat.User).Messages = make(chan interface{}, 5)
	repo.Users.Store(sessionId, user)

	//Test global message
	repo.broadcast(nil, chat.InitMessage, nil)
	repo.Users.Range(func(key, value interface{}) bool {
		user := value.(*chat.User)
		close(user.Messages)
		_, ok := <-user.Messages
		assert.Equal(t, true, ok, fmt.Sprintf("SessionId: %s. User doesn't receive global message", key.(string)))
		return true
	})

	close(repo.Connect)
	close(repo.Disconnect)
	close(repo.Actions)
}

var testCaseProcessAction = []struct {
	name           string
	action         *chat.Action
	isPrivate      bool
	ignoreReceiver bool
}{
	{
		name: "Test getting global messages",
		action: &chat.Action{
			Type: chat.InitGlobalMessages,
			Payload: &chat.InitGlobalMessagesPayload{
				Start:    0,
				Limit:    0,
				Receiver: &receiverId,
				MessageInfo: chat.MessageInfo{
					Author:    &receiverId,
					SessionID: receiverSessionId,
				},
			},
		},
		isPrivate: true,
	},
	{
		name: "Test sending private message",
		action: &chat.Action{
			Type: chat.InitMessage,
			Payload: &chat.InitMessagePayload{
				MessageInfo: chat.MessageInfo{
					Author:    &receiverId,
					SessionID: receiverSessionId,
				},
				Receiver: &receiverId,
				Data: chat.Data{
					Text: "text",
				},
			},
		},
		isPrivate: true,
	},
	{
		name: "Test sending global message",
		action: &chat.Action{
			Type: chat.InitMessage,
			Payload: &chat.InitMessagePayload{
				MessageInfo: chat.MessageInfo{
					Author:    &receiverId,
					SessionID: receiverSessionId,
				},
				Receiver: nil,
				Data: chat.Data{
					Text: "text",
				},
			},
		},
		isPrivate: false,
	},
	{
		name: "Test updating message",
		action: &chat.Action{
			Type: chat.InitUpdateMessage,
			Payload: &chat.InitUpdateMessagePayload{
				MessageInfo: chat.MessageInfo{
					SessionID: receiverSessionId,
					Author:    &receiverId,
				},
				Id: 1,
				Data: chat.Data{
					Text: "Text",
				},
			},
		},
		isPrivate: false,
	},
	{
		name: "Test printing",
		action: &chat.Action{
			Type: chat.InitPrinting,
			Payload: &chat.InitPrintingPayload{
				MessageInfo: chat.MessageInfo{
					Author:    &receiverId,
					SessionID: receiverSessionId,
				},
			},
		},
		isPrivate:      false,
		ignoreReceiver: true,
	},
	{
		name: "Test deleting message",
		action: &chat.Action{
			Type: chat.InitDeleteMessage,
			Payload: &chat.InitDeleteMessagePayload{
				MessageInfo: chat.MessageInfo{
					SessionID: receiverSessionId,
					Author:    &receiverId,
				},
				Id: 1,
			},
		},
		isPrivate: false,
	},
}

func TestProcessAction(t *testing.T) {
	for _, testCase := range testCaseProcessAction {
		t.Run(testCase.name, func(t *testing.T) {
			repo := repo()

			id := uint64(1)
			sessionId := uuid.NewV4().String()
			repo.Users.Store(sessionId, &chat.User{
				Id:        &id,
				SessionID: sessionId,
				Messages:  make(chan interface{}, 5),
			})
			repo.Sessions.Store(id, sessionId)

			repo.Users.Store(receiverSessionId, &chat.User{
				Id:        &receiverId,
				SessionID: receiverSessionId,
				Messages:  make(chan interface{}, 5),
			})
			repo.Sessions.Store(receiverId, receiverSessionId)

			repo.Actions <- testCase.action
			close(repo.Actions)
			repo.processAction()

			if testCase.isPrivate {
				sessionId, _ := repo.Sessions.Load(receiverId)
				user, _ := repo.Users.Load(sessionId)
				close(user.(*chat.User).Messages)
				_, ok := <-user.(*chat.User).Messages
				assert.Equal(t, true, ok, "Doesn't receive private messsage")
			} else {
				repo.Users.Range(func(key, value interface{}) bool {
					if testCase.ignoreReceiver && key.(string) == receiverSessionId {
						return true
					}
					user := value.(*chat.User)
					close(user.Messages)
					_, ok := <-user.Messages
					assert.Equal(t, true, ok, fmt.Sprintf("SessionId: %s. User doesn't receive global message", key.(string)))
					return true
				})
			}
		})
	}
}
