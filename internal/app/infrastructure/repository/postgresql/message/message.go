package message

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
	"github.com/jackc/pgx"
)

const (
	QueryCreateMessage = `INSERT INTO message (author, receiver, text)
		VALUES ($1, $2, $3)
		RETURNING id, created, edited`
)

// NewRepo creates new instance of chat repository
func NewRepo(conn *pgx.Conn) *Repo {
	return &Repo{
		conn: conn,
	}
}

type Repo struct {
	conn *pgx.Conn
}

// Create adds new massage
func (r *Repo) Create(message *chat.Message) (*chat.Message, error) {
	if err := r.conn.QueryRow(QueryCreateMessage, message.Author, message.Receiver, message.Data.Text).
		Scan(&message.ID, &message.Created, &message.Edited); err != nil {
		return nil, err
	}
	return message, nil
}
