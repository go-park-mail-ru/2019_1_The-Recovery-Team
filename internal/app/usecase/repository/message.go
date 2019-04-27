package repository

import "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"

type MessageRepo interface {
	Create(message *chat.Message) (*chat.Message, error)
}
