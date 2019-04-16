package game

import "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/domain/game"

type GameRepoMock struct{}

func (r *GameRepoMock) Run() {
	return
}

func (r *GameRepoMock) PlayersChan() chan *game.User {
	return make(chan *game.User)
}
