package repository

import "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/domain/profile"

type ProfileRepo interface {
	Get(id interface{}) (*profile.Profile, error)
	Create(data *profile.Create) (*profile.Created, error)
	Update(id interface{}, data *profile.UpdateInfo) error
	UpdateAvatar(id, avatarPath interface{}) error
	UpdatePassword(id interface{}, data *profile.UpdatePassword) error
	GetByEmail(email interface{}) (*profile.Profile, error)
	GetByNickname(nickname interface{}) (*profile.Profile, error)
	GetByEmailAndPassword(data *profile.Login) (*profile.Profile, error)
	List(limit, offset int64) ([]profile.Info, error)
	Count() (count int64, err error)
}
