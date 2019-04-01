package usecase

import (
	"sadislands/internal/usecase/repository"
	"time"
)

func NewSessionInteractor(repo repository.SessionRepo) *SessionInteractor {
	return &SessionInteractor{
		repo: repo,
	}
}

type SessionInteractor struct {
	repo repository.SessionRepo
}

func (i *SessionInteractor) Get(token string) (uint64, error) {
	return i.repo.Get(token)
}

func (i *SessionInteractor) Set(profileID uint64, expires time.Duration) (string, error) {
	return i.repo.Set(profileID, expires)
}

func (i *SessionInteractor) Delete(token string) error {
	return i.repo.Delete(token)
}
