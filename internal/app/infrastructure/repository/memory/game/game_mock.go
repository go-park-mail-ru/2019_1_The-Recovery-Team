package game

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/game"
)

type RepoMock struct{}

func (r *RepoMock) Run() {}

func (r *RepoMock) PlayersChan() chan *game.User {
	return make(chan *game.User)
}

func (r *RepoMock) UpdateRatingChan() chan *game.UpdateRating {
	return make(chan *game.UpdateRating)
}
