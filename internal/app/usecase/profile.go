package usecase

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase/repository"
)

func NewProfileInteractor(repo repository.ProfileRepo) *ProfileInteractor {
	return &ProfileInteractor{
		repo: repo,
	}
}

type ProfileInteractor struct {
	repo repository.ProfileRepo
}

func (i *ProfileInteractor) Get(id interface{}) (*profile.Profile, error) {
	return i.repo.Get(id)
}

func (i *ProfileInteractor) Create(data *profile.Create) (*profile.Created, error) {
	return i.repo.Create(data)
}

func (i *ProfileInteractor) Update(id interface{}, data *profile.UpdateInfo) error {
	return i.repo.Update(id, data)
}

func (i *ProfileInteractor) UpdateAvatar(id, avatarPath interface{}) error {
	return i.repo.UpdateAvatar(id, avatarPath)
}

func (i *ProfileInteractor) UpdatePassword(id interface{}, data *profile.UpdatePassword) error {
	return i.repo.UpdatePassword(id, data)
}

func (i *ProfileInteractor) GetByEmail(email interface{}) (*profile.Profile, error) {
	return i.repo.GetByEmail(email)
}

func (i *ProfileInteractor) GetByNickname(nickname interface{}) (*profile.Profile, error) {
	return i.repo.GetByNickname(nickname)
}

func (i *ProfileInteractor) GetByEmailAndPassword(data *profile.Login) (*profile.Profile, error) {
	return i.repo.GetByEmailAndPassword(data)
}

func (i *ProfileInteractor) List(limit, offset int64) ([]profile.Info, error) {
	return i.repo.List(limit, offset)
}

func (i *ProfileInteractor) Count() (count int64, err error) {
	return i.repo.Count()
}

func (i *ProfileInteractor) UpdateRating(winner, loser uint64) error {
	return i.repo.UpdateRating(winner, loser)
}
