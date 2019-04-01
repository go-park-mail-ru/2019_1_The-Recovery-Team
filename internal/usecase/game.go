package usecase

import (
	"sadislands/internal/domain/game"
	"sadislands/internal/usecase/repository"
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

func (i *GameInteractor) Players() chan *game.Player {
	return i.repo.PlayersChan()
}
