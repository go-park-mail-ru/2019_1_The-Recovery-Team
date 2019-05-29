package message

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/chat"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func repo() Repo {
	conn := postgresql.ConnMock{
		Tx: postgresql.TxMock{},
	}
	return Repo{
		conn: &conn,
	}
}

func TestNewChatRepo(t *testing.T) {
	conn := &pgx.ConnPool{}
	assert.NotEmpty(t, NewRepo(conn),
		"Doesn't create profile repository instance")
}

var testCaseCreate = []struct {
	name    string
	message chat.Message
	err     error
}{
	{
		name: "Test with incorrect author id",
		message: chat.Message{
			Author:   &postgresql.ForbiddenMessageAuthor,
			Receiver: nil,
			Data: chat.Data{
				Text: postgresql.MessageText,
			},
		},
		err: errors.New(DefaultError),
	},
	{
		name: "Test with correct data",
		message: chat.Message{
			Author:   &postgresql.MessageAuthor,
			Receiver: nil,
			Data: chat.Data{
				Text: postgresql.MessageText,
			},
		},
		err: nil,
	},
}

func TestList(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseCreate {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.Create(&testCase.message)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
		})
	}
}

var testCaseGetGlobal = []struct {
	name  string
	query chat.Query
	err   error
}{
	{
		name: "Test with incorrect limit",
		query: chat.Query{
			Limit: postgresql.ForbiddenMessageLimit,
		},
		err: errors.New(DefaultError),
	},
	{
		name: "Test with correct data",
		query: chat.Query{
			Limit: 1,
			Start: 1,
		},
		err: nil,
	},
}

func TestGetGlobal(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseGetGlobal {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.GetGlobal(&testCase.query)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
		})
	}
}

var testCaseUpdate = []struct {
	name    string
	message chat.Message
	err     error
}{
	{
		name: "Test wit incorrect message id",
		message: chat.Message{
			ID: postgresql.ForbiddenMessageId,
			Data: chat.Data{
				Text: postgresql.ForbiddenMessageText,
			},
		},
		err: errors.New(DefaultError),
	},
	{
		name: "Test without permission",
		message: chat.Message{
			ID:     postgresql.MessageId,
			Author: &postgresql.ForbiddenMessageAuthor,
			Data: chat.Data{
				Text: postgresql.ForbiddenMessageText,
			},
		},
		err: errors.New("permission denied"),
	},
	{
		name: "Test with incorrect message text",
		message: chat.Message{
			ID:     postgresql.MessageId,
			Author: &postgresql.MessageAuthor,
			Data: chat.Data{
				Text: postgresql.ForbiddenMessageText,
			},
		},
		err: errors.New(DefaultError),
	},
	{
		name: "Test with correct data",
		message: chat.Message{
			ID:     postgresql.MessageId,
			Author: &postgresql.MessageAuthor,
			Data: chat.Data{
				Text: postgresql.MessageText,
			},
		},
		err: nil,
	},
}

func TestUpdate(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseUpdate {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.Update(&testCase.message)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
		})
	}
}

var testCaseDelete = []struct {
	name    string
	message chat.Message
	err     error
}{
	{
		name: "Test wit incorrect message id",
		message: chat.Message{
			ID: postgresql.ForbiddenMessageId,
			Data: chat.Data{
				Text: postgresql.ForbiddenMessageText,
			},
		},
		err: errors.New(DefaultError),
	},
	{
		name: "Test without permission",
		message: chat.Message{
			ID:     postgresql.MessageId,
			Author: &postgresql.ForbiddenMessageAuthor,
			Data: chat.Data{
				Text: postgresql.ForbiddenMessageText,
			},
		},
		err: errors.New("permission denied"),
	},
	{
		name: "Test with correct data",
		message: chat.Message{
			ID:     postgresql.MessageId,
			Author: &postgresql.MessageAuthor,
			Data: chat.Data{
				Text: postgresql.MessageText,
			},
		},
		err: nil,
	},
}

func TestDelete(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseDelete {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.Delete(&testCase.message)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
		})
	}
}
