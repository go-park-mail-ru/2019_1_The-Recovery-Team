package repository

import "sadislands/internal/domain/game"

type GameRepo interface {
	Run()
	PlayersChan() chan *game.Player
}
