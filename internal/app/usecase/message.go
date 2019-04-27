package usecase

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase/repository"
)

func NewMessageInteractor(repo repository.MessageRepo) *MessageInteractor {
	return &MessageInteractor{
		repo: repo,
	}
}

type MessageInteractor struct {
	repo repository.MessageRepo
}

func (i *MessageInteractor) Create(message *chat.Message) (*chat.Message, error) {
	return i.repo.Create(message)
}

func (i *MessageInteractor) GetGlobal(data *chat.Query) (*[]chat.Message, error) {
	return i.repo.GetGlobal(data)
}

func (i *MessageInteractor) Update(message *chat.Message) (*chat.Message, error) {
	return i.repo.Update(message)
}

func (i *MessageInteractor) Delete(message *chat.Message) (*chat.Message, error) {
	return i.repo.Delete(message)
}
