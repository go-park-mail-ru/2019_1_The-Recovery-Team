package session

import (
	"errors"
	"time"
)

const (
	Authorized          = "AUTHORIZED"
	Unauthorized        = "UNAUTHORIZED"
	AuthorizedProfileId = 1
	DefaultProfileId    = 1
	IncorrectData       = "INCORRECT_DATA"
)

type SessionRepoMock struct{}

func (r *SessionRepoMock) Get(token string) (uint64, error) {
	switch token {
	case Authorized:
		{
			return AuthorizedProfileId, nil
		}
	default:
		{
			return 0, errors.New(Unauthorized)
		}
	}
}

func (r *SessionRepoMock) Set(profileID uint64, expires time.Duration) (string, error) {
	switch profileID {
	case DefaultProfileId:
		{
			return Authorized, nil
		}
	default:
		{
			return "", errors.New(IncorrectData)
		}
	}
}

func (r *SessionRepoMock) Delete(token string) error {
	if token != Authorized {
		return errors.New(Unauthorized)
	}

	return nil
}
