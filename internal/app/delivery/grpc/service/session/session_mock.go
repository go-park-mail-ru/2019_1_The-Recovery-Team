package session

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/redis/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	"google.golang.org/grpc"
)

type ClientMock struct {
	service *Service
}

func NewClientMock() SessionClient {
	return &ClientMock{
		service: NewService(usecase.NewSessionInteractor(&session.RepoMock{})),
	}
}

func (c *ClientMock) Get(ctx context.Context, in *SessionId, opts ...grpc.CallOption) (*ProfileId, error) {
	return c.service.Get(ctx, in)
}

func (c *ClientMock) Set(ctx context.Context, in *Create, opts ...grpc.CallOption) (*SessionId, error) {
	return c.service.Set(ctx, in)
}

func (c *ClientMock) Delete(ctx context.Context, in *SessionId, opts ...grpc.CallOption) (*Nothing, error) {
	return c.service.Delete(ctx, in)
}
