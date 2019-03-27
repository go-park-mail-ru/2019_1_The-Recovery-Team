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

func (s *SessionInteractor) Get(token string) (uint64, error) {
	return s.repo.Get(token)
}

func (s *SessionInteractor) Set(profileID uint64, expires time.Duration) (string, error) {
	return s.repo.Set(profileID, expires)
}

func (s *SessionInteractor) Delete(token string) error {
	return s.repo.Delete(token)
}
