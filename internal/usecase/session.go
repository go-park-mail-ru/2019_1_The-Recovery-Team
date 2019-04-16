package usecase

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase/repository"
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
