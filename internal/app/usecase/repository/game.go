package repository

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/game"
)

type GameRepo interface {
	Run()
	PlayersChan() chan *game.User
}
