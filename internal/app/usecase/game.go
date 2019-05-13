package usecase

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase/repository"
)

func NewGameInteractor(repo repository.GameRepo) *GameInteractor {
	return &GameInteractor{
		repo: repo,
	}
}

type GameInteractor struct {
	repo repository.GameRepo
}

func (i *GameInteractor) Run() {
	i.repo.Run()
}

func (i *GameInteractor) Players() chan *game.User {
	return i.repo.PlayersChan()
}

func (i *GameInteractor) UpdateRatingChan() chan *game.UpdateRating {
	return i.repo.UpdateRatingChan()
}
