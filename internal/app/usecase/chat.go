package usecase

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase/repository"
)

func NewChatInteractor(repo repository.ChatRepo) *ChatInteractor {
	return &ChatInteractor{
		repo: repo,
	}
}

type ChatInteractor struct {
	repo repository.ChatRepo
}

func (i *ChatInteractor) Run() {
	i.repo.Run()
}

func (i *ChatInteractor) Connection() chan *chat.User {
	return i.repo.Connection()
}
