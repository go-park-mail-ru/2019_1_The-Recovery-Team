package session

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
)

func NewService(interactor *usecase.SessionInteractor) *Service {
	return &Service{
		interactor: interactor,
	}
}

type Service struct {
	interactor *usecase.SessionInteractor
}

func (s *Service) Get(ctx context.Context, in *SessionId) (*ProfileId, error) {
	profileId, err := s.interactor.Get(in.Id)
	id := &ProfileId{Id: profileId}
	return id, err
}

func (s *Service) Set(ctx context.Context, in *Create) (*SessionId, error) {
	duration := time.Duration(in.Expires.Seconds)
	id, err := s.interactor.Set(in.ProfileId.Id, duration)
	session := &SessionId{Id: id}
	return session, err
}

func (s *Service) Delete(ctx context.Context, in *SessionId) (*Nothing, error) {
	err := s.interactor.Delete(in.Id)
	return &Nothing{}, err
}
