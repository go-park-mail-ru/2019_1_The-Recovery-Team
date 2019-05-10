package session

import (
	"context"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
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
	if err != nil {
		return &ProfileId{}, err
	}

	return id, nil
}

func (s *Service) Set(ctx context.Context, in *Create) (*SessionId, error) {
	duration, err := time.ParseDuration(strconv.Itoa(int((*in.Expires).Seconds)) + "s")
	if err != nil {
		return &SessionId{}, err
	}

	id, err := s.interactor.Set(in.ProfileId.Id, duration)
	if err != nil {
		return &SessionId{}, err
	}
	session := &SessionId{Id: id}
	return session, nil
}

func (s *Service) Delete(ctx context.Context, in *SessionId) (*Nothing, error) {
	err := s.interactor.Delete(in.Id)
	return &Nothing{}, err
}
