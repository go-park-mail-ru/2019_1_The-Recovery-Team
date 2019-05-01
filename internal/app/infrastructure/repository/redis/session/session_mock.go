package session

import (
	"errors"
	"time"
)

const (
	Authorized                 = "AUTHORIZED"
	Unauthorized               = "UNAUTHORIZED"
	AuthorizedProfileId        = 1
	DefaultProfileId    uint64 = 1
	IncorrectData              = "INCORRECT_DATA"
)

type RepoMock struct{}

func (r *RepoMock) Get(token string) (uint64, error) {
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

func (r *RepoMock) Set(profileID uint64, expires time.Duration) (string, error) {
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

func (r *RepoMock) Delete(token string) error {
	if token != Authorized {
		return errors.New(Unauthorized)
	}

	return nil
}
