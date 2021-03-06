package profile

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
)

func NewService(interactor *usecase.ProfileInteractor) *Service {
	return &Service{
		interactor: interactor,
	}
}

type Service struct {
	interactor *usecase.ProfileInteractor
}

func (s *Service) Get(ctx context.Context, in *GetRequest) (*GetResponse, error) {
	prof, err := s.interactor.Get(in.Id)
	if err != nil {
		return &GetResponse{}, err
	}

	r := &GetResponse{
		Info: &Info{
			Id:       prof.ID,
			Nickname: prof.Nickname,
			Avatar:   prof.Avatar,
			Oauth:    prof.Oauth,
			OauthId:  prof.OauthId,
			Score: &Score{
				Position: prof.Score.Position,
				Record:   prof.Score.Record,
				Win:      prof.Score.Win,
				Loss:     prof.Score.Loss,
			},
		},
		Email: prof.Email,
	}
	return r, nil
}

func (s *Service) Create(ctx context.Context, in *CreateRequest) (*CreateResponse, error) {
	create := &profile.Create{
		Email:    in.Email,
		Nickname: in.Nickname,
		Password: in.Password,
	}
	created, err := s.interactor.Create(create)
	if err != nil {
		return &CreateResponse{}, err
	}

	r := &CreateResponse{
		Id:       created.ID,
		Email:    created.Email,
		Nickname: created.Nickname,
		Avatar:   created.Avatar,
		Score: &Score{
			Position: created.Score.Position,
			Record:   created.Score.Record,
			Win:      created.Score.Win,
			Loss:     created.Score.Loss,
		},
	}
	return r, nil
}

func (s *Service) Update(ctx context.Context, in *UpdateRequest) (*Nothing, error) {
	update := &profile.UpdateInfo{
		Email:    in.Email,
		Nickname: in.Nickname,
	}
	err := s.interactor.Update(in.Id, update)
	return &Nothing{}, err
}

func (s *Service) UpdateAvatar(ctx context.Context, in *UpdateAvatarRequest) (*Nothing, error) {
	err := s.interactor.UpdateAvatar(in.Id, in.Avatar)
	return &Nothing{}, err
}

func (s *Service) UpdatePassword(ctx context.Context, in *UpdatePasswordRequest) (*Nothing, error) {
	update := &profile.UpdatePassword{
		Password:    in.Password,
		PasswordOld: in.PasswordOld,
	}
	err := s.interactor.UpdatePassword(in.Id, update)
	return &Nothing{}, err
}

func (s *Service) GetByEmail(ctx context.Context, in *GetByEmailRequest) (*GetResponse, error) {
	prof, err := s.interactor.GetByEmail(in.Email)
	if err != nil {
		return &GetResponse{}, err
	}

	r := &GetResponse{
		Info: &Info{
			Id:       prof.ID,
			Nickname: prof.Nickname,
			Avatar:   prof.Avatar,
			Score: &Score{
				Record: prof.Score.Record,
				Win:    prof.Score.Win,
				Loss:   prof.Score.Loss,
			},
		},
		Email: prof.Email,
	}
	return r, nil
}

func (s *Service) GetByNickname(ctx context.Context, in *GetByNicknameRequest) (*GetResponse, error) {
	prof, err := s.interactor.GetByNickname(in.Nickname)
	if err != nil {
		return &GetResponse{}, err
	}

	r := &GetResponse{
		Info: &Info{
			Id:       prof.ID,
			Nickname: prof.Nickname,
			Avatar:   prof.Avatar,
			Score: &Score{
				Record: prof.Score.Record,
				Win:    prof.Score.Win,
				Loss:   prof.Score.Loss,
			},
		},
		Email: prof.Email,
	}
	return r, nil
}

func (s *Service) GetByEmailAndPassword(ctx context.Context, in *GetByEmailAndPasswordRequest) (*GetResponse, error) {
	login := &profile.Login{
		Email:    in.Email,
		Password: in.Password,
	}
	prof, err := s.interactor.GetByEmailAndPassword(login)
	if err != nil {
		return &GetResponse{}, err
	}

	r := &GetResponse{
		Info: &Info{
			Id:       prof.ID,
			Nickname: prof.Nickname,
			Avatar:   prof.Avatar,
			Score: &Score{
				Record:   prof.Score.Record,
				Win:      prof.Score.Win,
				Loss:     prof.Score.Loss,
				Position: prof.Score.Position,
			},
		},
		Email: prof.Email,
	}
	return r, nil
}

func (s *Service) List(ctx context.Context, in *ListRequest) (*ListResponse, error) {
	list, err := s.interactor.List(in.Limit, in.Offset)
	if err != nil {
		return &ListResponse{}, err
	}

	r := &ListResponse{
		List: make([]*Info, 0, len(list)),
	}
	for _, info := range list {
		r.List = append(r.List,
			&Info{
				Id:       info.ID,
				Nickname: info.Nickname,
				Avatar:   info.Avatar,
				Score: &Score{
					Record: info.Score.Record,
					Win:    info.Score.Win,
					Loss:   info.Score.Loss,
				},
			})
	}
	return r, nil
}

func (s *Service) Count(ctx context.Context, in *Nothing) (*CountResponse, error) {
	count, err := s.interactor.Count()
	if err != nil {
		return &CountResponse{}, err
	}

	r := &CountResponse{
		Count: count,
	}
	return r, nil
}

func (s *Service) UpdateRating(ctx context.Context, in *UpdateRatingRequest) (*Nothing, error) {
	err := s.interactor.UpdateRating(in.Winner, in.Loser)
	return &Nothing{}, err
}

func (s *Service) PutProfileOauth(ctx context.Context, in *PutProfileOauthRequest) (*ProfileId, error) {
	p, err := s.interactor.PutProfileOauth(in.Id, in.Token)
	if err != nil {
		return &ProfileId{}, err
	}

	return &ProfileId{Id: p.Id}, nil
}

func (s *Service) CreateProfileOauth(ctx context.Context, in *CreateProfileOauthRequest) (*ProfileId, error) {
	c := &profile.CreateOauth{
		UserId: in.UserId,
		Token:  in.Token,
		Avatar: profile.Avatar{
			Path: in.Avatar,
		},
		Oauth: in.Oauth,
	}
	p, err := s.interactor.CreateProfileOauth(c)
	if err != nil {
		return &ProfileId{}, err
	}

	return &ProfileId{Id: p.Id}, nil
}
