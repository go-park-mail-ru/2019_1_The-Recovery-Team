package profile

import (
	"errors"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/profile"

	"github.com/jackc/pgx"
)

const (
	DefaultProfileId    = 1
	DefaultProfileIdStr = "1"

	ExistingProfileId       = 2
	ExistingProfileIdStr    = "2"
	ExistingProfileEmail    = "test@mail.ru"
	ExistingProfileNickname = "test"
	ExistingProfilePassword = "1234"

	NotExistingProfileIdStr    = "101"
	NotExistingProfileId       = 101
	NotExistingProfileEmail    = "notExist@mail.ru"
	NotExistingProfileNickname = "notExist"
	NotExistingProfilePassword = "notExist"

	IncorrectProfileEmail = "email"

	InvalidProfilePassword = "bad"

	CreatedProfileEmail    = "new@mail.ru"
	CreatedProfileNickname = "new"

	ForbiddenProfileId       = 100
	ForbiddenProfileIdStr    = "100"
	ForbiddenProfileEmail    = "forbidden@mail.ru"
	ForbiddenProfileNickname = "forbidden"
	ForbiddenProfileAvatar   = "forbidden"
	ForbiddenProfilePassword = "password"
	ForbiddenLimit           = -1
	ForbiddenLimitStr        = "-1"

	DefaultError = "error"
	DefaultCount = 1
)

type RepoMock struct{}

func (r *RepoMock) Get(id interface{}) (*profile.Profile, error) {
	profileId, _ := id.(uint64)
	if id.(uint64) == ForbiddenProfileId {
		return nil, errors.New(DefaultError)
	}

	switch profileId {
	case ExistingProfileId, DefaultProfileId:
		{
			prof := &profile.Profile{
				Info: profile.Info{
					ID: profileId,
				},
			}
			return prof, nil
		}
	default:
		{
			return nil, pgx.ErrNoRows
		}
	}
}

func (r *RepoMock) Create(data *profile.Create) (*profile.Created, error) {
	if data.Email == ForbiddenProfileEmail {
		return nil, errors.New(DefaultError)
	}

	if data.Nickname == ExistingProfileNickname {
		return nil, errors.New(NicknameAlreadyExists)
	}

	if data.Email == ExistingProfileEmail {
		return nil, errors.New(EmailAlreadyExists)
	}

	if matches, _ := postgresql.VerifyPassword(IncorrectProfilePassword, data.Password); matches {
		created := &profile.Created{
			ID:       ForbiddenProfileId,
			Email:    CreatedProfileEmail,
			Nickname: CreatedProfileNickname,
		}
		return created, nil
	}

	created := &profile.Created{
		ID:       DefaultProfileId,
		Email:    CreatedProfileEmail,
		Nickname: CreatedProfileNickname,
	}
	return created, nil
}

func (r *RepoMock) Update(id interface{}, data *profile.UpdateInfo) error {
	if data.Email == ForbiddenProfileEmail {
		return errors.New(DefaultError)
	}

	if data.Nickname == ExistingProfileNickname {
		return errors.New(NicknameAlreadyExists)
	}

	if data.Email == ExistingProfileEmail {
		return errors.New(EmailAlreadyExists)
	}

	return nil
}

func (r *RepoMock) UpdateAvatar(id, avatarPath interface{}) error {
	if avatarPath.(string) == ForbiddenProfileAvatar {
		return errors.New(DefaultError)
	}

	return nil
}

func (r *RepoMock) UpdatePassword(id interface{}, data *profile.UpdatePassword) error {
	if data.Password == ForbiddenProfilePassword || data.PasswordOld == ForbiddenProfilePassword {
		return errors.New(DefaultError)
	}

	if data.PasswordOld != ExistingProfilePassword {
		return errors.New(IncorrectProfilePassword)
	}

	return nil
}

func (r *RepoMock) GetByEmail(email interface{}) (*profile.Profile, error) {
	profileEmail := email.(string)

	if profileEmail == ForbiddenProfileEmail {
		return nil, errors.New(DefaultError)
	}

	if profileEmail != ExistingProfileEmail {
		return nil, pgx.ErrNoRows
	}

	received := &profile.Profile{
		Email: profileEmail,
	}

	return received, nil
}

func (r *RepoMock) GetByNickname(nickname interface{}) (*profile.Profile, error) {
	profileNickname := nickname.(string)

	if profileNickname == ForbiddenProfileNickname {
		return nil, errors.New(DefaultError)
	}

	if profileNickname != ExistingProfileNickname {
		return nil, pgx.ErrNoRows
	}

	received := &profile.Profile{
		Info: profile.Info{
			Nickname: profileNickname,
		},
	}

	return received, nil
}

func (r *RepoMock) GetByEmailAndPassword(data *profile.Login) (*profile.Profile, error) {
	if data.Email != ExistingProfileEmail {
		return nil, errors.New(DefaultError)
	}

	if data.Password != ExistingProfilePassword {
		return nil, pgx.ErrNoRows
	}

	if data.Email == ForbiddenProfileEmail {
		received := &profile.Profile{
			Info: profile.Info{
				ID: ExistingProfileId,
			},
		}
		return received, nil
	}

	received := &profile.Profile{
		Info: profile.Info{
			ID: DefaultProfileId,
		},
	}
	return received, nil
}

func (r *RepoMock) List(limit, offset int64) ([]profile.Info, error) {
	if limit == ForbiddenLimit {
		return nil, errors.New(DefaultError)
	}

	profiles := []profile.Info{
		{ID: DefaultProfileId},
	}
	return profiles, nil
}

func (r *RepoMock) Count() (count int64, err error) {
	count = DefaultCount
	err = nil
	return
}
