package repository

import "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/domain/profile"

type ProfileRepo interface {
	GetProfile(id interface{}) (*profile.Profile, error)
	CreateProfile(data *profile.Create) (*profile.Created, error)
	UpdateProfile(id interface{}, data *profile.UpdateInfo) error
	UpdateProfileAvatar(id, avatarPath interface{}) error
	UpdateProfilePassword(id interface{}, data *profile.UpdatePassword) error
	GetProfileByEmail(email interface{}) (*profile.Profile, error)
	GetProfileByNickname(nickname interface{}) (*profile.Profile, error)
	GetProfileByEmailWithPassword(data *profile.Login) (*profile.Profile, error)
	GetProfiles(limit, offset int64) ([]profile.Info, error)
	GetProfileCount() (count int64, err error)
}
