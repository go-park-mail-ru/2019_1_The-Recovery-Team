package repository

import "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"

type ChatRepo interface {
	Run()
	Connection() chan *chat.User
}
