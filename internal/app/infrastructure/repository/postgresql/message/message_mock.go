package message

import (
	"errors"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
)

const (
	CorrectMessageText  = "text"
	FobiddenMessageText = "forbidden"
	CorrectQueryLimit   = 1
	ForbiddenQueryLimit = 100

	ExistingMessageId         = 1
	ExistingAuthorId   uint64 = 1
	ExistingReceiverId uint64 = 2

	DefaultError = "error"
)

type RepoMock struct{}

func (r *RepoMock) Create(message *chat.Message) (*chat.Message, error) {
	if message.Data.Text != CorrectMessageText {
		return nil, errors.New(DefaultError)
	}

	message.ID = ExistingMessageId
	message.Created = time.Now()
	return message, nil
}

func (r *RepoMock) GetGlobal(data *chat.Query) (*[]chat.Message, error) {
	if data.Limit == ForbiddenQueryLimit {
		return nil, errors.New(DefaultError)
	}

	messages := make([]chat.Message, 0, 1)
	messages = append(messages, chat.Message{
		ID:       ExistingMessageId,
		Author:   data.Author,
		Receiver: data.Receiver,
		Created:  time.Now(),
		Data: chat.Data{
			Text: CorrectMessageText,
		},
	})

	return &messages, nil
}

func (r *RepoMock) Update(message *chat.Message) (*chat.Message, error) {
	return message, nil
}

func (r *RepoMock) Delete(message *chat.Message) (*chat.Message, error) {
	return message, nil
}
