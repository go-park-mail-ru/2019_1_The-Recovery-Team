package postgresql

import (
	"errors"
	"reflect"
	"time"

	"github.com/jackc/pgx"
)

const (
	QueryGetMesssageAuthor = `SELECT author FROM message
		WHERE id = $1`

	QueryGetMesssage = `SELECT id, author, receiver, created, edited, text FROM message
		WHERE id = $1`

	QueryCreateMessage = `INSERT INTO message (author, receiver, text)
		VALUES ($1, $2, $3)
		RETURNING id, created, edited`

	QueryGetGlobalMessagesFrom = `SELECT id, author, receiver, created, edited, text FROM message
		WHERE id < $1
		ORDER BY created DESC
		LIMIT $2`

	QueryGetGlobalMessages = `SELECT id, author, receiver, created, edited, text FROM message
		WHERE id > $1
		ORDER BY created DESC
		LIMIT $2`

	QueryUpdateMessage = `UPDATE message SET text = $1, edited = true WHERE id = $2
		RETURNING id, author, receiver, created, edited`

	QueryDeleteMessage = `DELETE FROM message WHERE id = $1`

	QueryProfileById = `SELECT id, nickname, email, avatar, record, win, loss 
    FROM profile 
    WHERE id = $1`

	QueryCreateProfile = `INSERT INTO profile (email, nickname, password) 
	VALUES ($1, $2, $3) 
	RETURNING id, email, nickname`

	QueryUpdateProfile = `UPDATE profile
	SET email = (CASE WHEN $1 = '' THEN email ELSE $1 END),
	nickname = (CASE WHEN $2 = '' THEN nickname ELSE $2 END)
	WHERE id = $3`

	QueryUpdateProfileAvatar = `UPDATE profile
	SET avatar = $1
	WHERE id = $2`

	QueryUpdateProfilePassword = `UPDATE profile
	SET password = $1
	WHERE id = $2`

	QueryProfileByEmail = `SELECT id, email, nickname
	FROM profile
	WHERE email = $1`

	QueryProfileByNickname = `SELECT id, email, nickname
	FROM profile
	WHERE nickname = $1`

	QueryProfileByIdWithPassword = `SELECT password
	FROM profile
	WHERE id = $1`

	QueryProfileByEmailWithPassword = `SELECT id, email, nickname, password, avatar, record, win, loss  
	FROM profile 
	WHERE email = $1`

	QueryProfilesWithLimitAndOffset = `SELECT id, nickname, avatar, record, win, loss 
	FROM profile 
	ORDER BY record LIMIT $1 OFFSET $2`

	QueryProfileCount = `SELECT COUNT(*) 
	FROM profile`

	ProfileEmailKey    = "profile_email_key"
	ProfileNicknameKey = "profile_nickname_key"
)

var (
	ProfileId       uint64 = 1
	ProfileEmail           = "Email"
	ProfileNickname        = "Nickname"
	ProfileAvatar          = "Avatar"
	ProfileScore    int64  = 0
	ProfilePassword        = "password"
	ProfileCount    int64  = 1
	MessageAuthor   uint64 = 1
	MessageReceiver uint64 = 1
	MessageText            = "text"
	MessageId       uint64 = 1
	MessageCreated         = time.Now()
	MessageEdited          = false

	ConflictProfileEmail    = "conflict"
	ConflictProfileNickname = "conflict"

	ForbiddenProfileId     uint64 = 101
	ForbiddenEmail                = "forbidden"
	ForbiddenLimit         int64  = 101
	ForbiddenMessageAuthor uint64 = 101
	ForbiddenMessageLimit  int    = 101
	ForbiddenMessageId     uint64 = 101
	ForbiddenMessageText          = "forbidden"

	DefaultError = "error"
)

type ConnMock struct {
	Tx TxMock
}

func (c *ConnMock) Close() {}

func (c *ConnMock) Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error) {
	return "", nil
}

func (c *ConnMock) Query(sql string, args ...interface{}) (Rows, error) {
	switch sql {
	case QueryProfilesWithLimitAndOffset:
		{
			if args[0].(int64) == ForbiddenLimit {
				return &RowsMock{}, errors.New(DefaultError)
			}

			row := RowMock{
				values: []interface{}{
					&ProfileId,
					&ProfileNickname,
					&ProfileAvatar,
					&ProfileScore,
					&ProfileScore,
					&ProfileScore,
				},
			}

			rows := RowsMock{
				values: []RowMock{
					row,
				},
			}
			return &rows, nil
		}
	case QueryGetGlobalMessages, QueryGetGlobalMessagesFrom:
		{
			if args[1].(int) == ForbiddenMessageLimit {
				return nil, errors.New(DefaultError)
			}

			row := RowMock{
				values: []interface{}{
					&MessageId,
					&MessageAuthor,
					&MessageReceiver,
					&MessageCreated,
					&MessageEdited,
					&MessageText,
				},
			}

			rows := RowsMock{
				values: []RowMock{
					row,
				},
			}
			return &rows, nil
		}
	}
	return nil, nil
}

func (c *ConnMock) QueryRow(sql string, args ...interface{}) Row {
	switch sql {
	case QueryProfileById:
		{
			if args[0].(uint64) == ForbiddenProfileId {
				return &RowMock{
					err: errors.New(DefaultError),
				}
			}

			row := RowMock{
				values: []interface{}{
					&ProfileId,
					&ProfileNickname,
					&ProfileEmail,
					&ProfileAvatar,
					&ProfileScore,
					&ProfileScore,
					&ProfileScore,
				},
			}
			return &row
		}
	case QueryCreateProfile:
		{
			if args[0].(uint64) == ForbiddenProfileId {
				return &RowMock{
					err: errors.New(DefaultError),
				}
			}

			row := RowMock{
				values: []interface{}{
					&ProfileId,
					&ProfileNickname,
					&ProfileEmail,
					&ProfileAvatar,
					&ProfileScore,
					&ProfileScore,
					&ProfileScore,
				},
			}
			return &row
		}
	case QueryProfileByIdWithPassword:
		{
			if args[0].(uint64) == ForbiddenProfileId {
				row := RowMock{
					err: errors.New(DefaultError),
				}
				return &row
			}

			password, err := HashAndSalt(ProfilePassword)
			if err != nil {
				panic(err)
			}

			row := RowMock{
				values: []interface{}{
					&password,
				},
			}
			return &row
		}
	case QueryProfileByEmail, QueryProfileByNickname:
		{
			return &RowMock{
				values: []interface{}{
					&ProfileId,
					&ProfileNickname,
					&ProfileEmail,
				},
			}
		}
	case QueryProfileByEmailWithPassword:
		{
			if args[0].(string) == ForbiddenEmail {
				return &RowMock{
					err: errors.New(DefaultError),
				}
			}

			password, err := HashAndSalt(ProfilePassword)
			if err != nil {
				panic(err)
			}
			return &RowMock{
				values: []interface{}{
					&ProfileId,
					&ProfileNickname,
					&ProfileEmail,
					&password,
					&ProfileAvatar,
					&ProfileScore,
					&ProfileScore,
					&ProfileScore,
				},
			}
		}
	case QueryProfileCount:
		{
			return &RowMock{
				values: []interface{}{
					&ProfileCount,
				},
			}
		}
	case QueryCreateMessage:
		{
			if *args[0].(*uint64) == ForbiddenMessageAuthor {
				return &RowMock{
					err: errors.New(DefaultError),
				}
			}

			row := RowMock{
				values: []interface{}{
					&MessageId,
					&MessageCreated,
					&MessageEdited,
				},
			}
			return &row
		}
	}
	return nil
}

func (c *ConnMock) Begin() (Tx, error) {
	return &c.Tx, nil
}

type TxMock struct{}

func (t *TxMock) Commit() error {
	return nil
}

func (t *TxMock) Rollback() error {
	return nil
}

func (t *TxMock) Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error) {
	switch sql {
	case QueryUpdateProfile:
		{
			if arguments[0].(string) == ConflictProfileEmail {
				err := pgx.PgError{
					ConstraintName: ProfileEmailKey,
				}
				return "", err
			}

			if arguments[1].(string) == ConflictProfileNickname {
				err := pgx.PgError{
					ConstraintName: ProfileNicknameKey,
				}
				return "", err
			}

			if arguments[2].(uint64) == ForbiddenProfileId {
				return "", errors.New(DefaultError)
			}

			return "", nil
		}
	case QueryUpdateProfileAvatar, QueryUpdateProfilePassword:
		{
			return "", nil
		}
	}
	return "", errors.New(DefaultError)
}

func (t *TxMock) QueryRow(sql string, args ...interface{}) Row {
	switch sql {
	case QueryCreateProfile:
		{
			if args[0].(string) == ConflictProfileEmail {
				row := RowMock{
					err: pgx.PgError{
						ConstraintName: ProfileEmailKey,
					},
				}
				return &row
			}

			if args[1].(string) == ConflictProfileNickname {
				row := RowMock{
					err: pgx.PgError{
						ConstraintName: ProfileNicknameKey,
					},
				}
				return &row
			}

			row := RowMock{
				values: []interface{}{
					&ProfileId,
					&ProfileNickname,
					&ProfileEmail,
				},
			}
			return &row
		}
	case QueryGetMesssageAuthor:
		{
			if args[0].(uint64) == ForbiddenMessageId {
				return &RowMock{
					err: errors.New(DefaultError),
				}
			}

			row := RowMock{
				values: []interface{}{
					&MessageAuthor,
				},
			}
			return &row
		}
	case QueryUpdateMessage:
		{
			if args[0].(string) == ForbiddenMessageText {
				return &RowMock{
					err: errors.New(DefaultError),
				}
			}

			row := RowMock{
				values: []interface{}{
					&MessageId,
					&MessageAuthor,
					&MessageReceiver,
					&MessageCreated,
					&MessageEdited,
				},
			}
			return &row
		}
	case QueryGetMesssage:
		{
			if args[0].(uint64) == ForbiddenMessageId {
				return &RowMock{
					err: errors.New(DefaultError),
				}
			}

			row := RowMock{
				values: []interface{}{
					&MessageId,
					&MessageAuthor,
					&MessageReceiver,
					&MessageCreated,
					&MessageEdited,
					&MessageText,
				},
			}
			return &row
		}
	case QueryDeleteMessage:
		{
			return &RowMock{}
		}
	}
	return nil
}

type RowMock struct {
	values []interface{}
	err    error
}

func (r *RowMock) Scan(dest ...interface{}) (err error) {
	if r.err != nil {
		return r.err
	}

	for i := 0; i < len(dest); i++ {
		switch reflect.TypeOf(dest[i]).String() {
		case "*uint64":
			{
				*dest[i].(*uint64) = *r.values[i].(*uint64)
			}
		case "*int64":
			{
				*dest[i].(*int64) = *r.values[i].(*int64)
			}
		case "*string":
			{
				*dest[i].(*string) = *r.values[i].(*string)
			}
		case "*time.Time":
			{
				*dest[i].(*time.Time) = *r.values[i].(*time.Time)
			}
		}
	}
	return
}

type RowsMock struct {
	values  []RowMock
	current int
}

func (r *RowsMock) Next() bool {
	return len(r.values) == r.current
}

func (r *RowsMock) Scan(dest ...interface{}) (err error) {
	err = r.values[r.current].Scan(dest...)
	r.current++
	return
}
