package chat

import "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"

type RepoMock struct{}

func (r *RepoMock) Run() {}

func (r *RepoMock) Connection() chan *chat.User {
	return make(chan *chat.User, 1)
}
