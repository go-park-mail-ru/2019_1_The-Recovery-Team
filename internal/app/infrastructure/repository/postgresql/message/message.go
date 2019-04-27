package message

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
	"github.com/jackc/pgx"
)

const (
	QueryCreateMessage = `INSERT INTO message (author, receiver, text)
		VALUES ($1, $2, $3)
		RETURNING id, created, edited`

	QueryGetMessages = `SELECT * FROM message
		WHERE ((author = $1 AND receiver = $2) OR (author = $2 AND receiver = $1)) AND id > $3
		ORDER BY created DESC
		LIMIT $4`

	QueryGetGlobalMessagesFrom = `SELECT id, author, receiver, created, edited, text FROM message
		WHERE id < $1
		ORDER BY created DESC
		LIMIT $2`

	QueryGetGlobalMessages = `SELECT id, author, receiver, created, edited, text FROM message
		WHERE id > $1
		ORDER BY created DESC
		LIMIT $2`

	QueryUpdateMessage = `UPDATE message SET text = $1, isEdited = true WHERE id = $2
		RETURNING id, author, receiver, created, edited`
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

// GetGlobal gets messages
func (r *Repo) GetGlobal(data *chat.Query) (*[]chat.Message, error) {
	messages := make([]chat.Message, 0, 10)

	var query string
	if data.Start == 0 {
		query = QueryGetGlobalMessages
	} else {
		query = QueryGetGlobalMessagesFrom
	}
	rows, err := r.conn.Query(query, data.Start, data.Limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		message := chat.Message{
			Data: chat.Data{},
		}
		if err = rows.Scan(&message.ID, &message.Author, &message.Receiver, &message.Created, &message.Edited, &message.Data.Text); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return &messages, nil
}

// Update text of message
func (r *Repo) Update(message *chat.Message) (*chat.Message, error) {
	if err := r.conn.QueryRow(QueryUpdateMessage, message.Data.Text, message.ID).
		Scan(&message.ID, &message.Author, &message.Receiver, &message.Created, &message.Edited); err != nil {
		return nil, err
	}
	return message, nil
}
