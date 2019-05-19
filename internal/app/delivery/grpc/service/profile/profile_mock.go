package profile

import (
	"context"

	repo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"google.golang.org/grpc"
)

type ClientMock struct {
	service *Service
}

func NewClientMock() ProfileClient {
	return &ClientMock{
		service: NewService(usecase.NewProfileInteractor(&repo.RepoMock{})),
	}
}

func (c *ClientMock) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	return c.service.Get(ctx, in)
}

func (c *ClientMock) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	return c.service.Create(ctx, in)
}

func (c *ClientMock) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*Nothing, error) {
	return c.service.Update(ctx, in)
}

func (c *ClientMock) UpdateAvatar(ctx context.Context, in *UpdateAvatarRequest, opts ...grpc.CallOption) (*Nothing, error) {
	return c.service.UpdateAvatar(ctx, in)
}

func (c *ClientMock) UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*Nothing, error) {
	return c.service.UpdatePassword(ctx, in)
}

func (c *ClientMock) GetByEmail(ctx context.Context, in *GetByEmailRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	return c.service.GetByEmail(ctx, in)
}

func (c *ClientMock) GetByNickname(ctx context.Context, in *GetByNicknameRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	return c.service.GetByNickname(ctx, in)
}

func (c *ClientMock) GetByEmailAndPassword(ctx context.Context, in *GetByEmailAndPasswordRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	return c.service.GetByEmailAndPassword(ctx, in)
}

func (c *ClientMock) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	return c.service.List(ctx, in)
}

func (c *ClientMock) Count(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*CountResponse, error) {
	return c.service.Count(ctx, in)
}

func (c *ClientMock) UpdateRating(ctx context.Context, in *UpdateRatingRequest, opts ...grpc.CallOption) (*Nothing, error) {
	return c.service.UpdateRating(ctx, in)
}

func (c *ClientMock) PutProfileOauth(ctx context.Context, in *PutProfileOauthRequest, opts ...grpc.CallOption) (*ProfileId, error) {
	return c.service.PutProfileOauth(ctx, in)
}

func (c *ClientMock) CreateProfileOauth(ctx context.Context, in *CreateProfileOauthRequest, opts ...grpc.CallOption) (*ProfileId, error) {
	return c.service.CreateProfileOauth(ctx, in)
}
