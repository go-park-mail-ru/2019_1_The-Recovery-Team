package game

import "sadislands/internal/domain/game"

type GameRepoMock struct{}

func (r *GameRepoMock) Run() {
	return
}

func (r *GameRepoMock) PlayersChan() chan *game.User {
	return make(chan *game.User)
}
