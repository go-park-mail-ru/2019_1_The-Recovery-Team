package usecase

import (
	"sadislands/internal/domain/profile"
	"sadislands/internal/usecase/repository"
)

func NewProfileInteractor(repo repository.ProfileRepo) *ProfileInteractor {
	return &ProfileInteractor{
		repo: repo,
	}
}

type ProfileInteractor struct {
	repo repository.ProfileRepo
}

func (i *ProfileInteractor) GetProfile(id interface{}) (*profile.Profile, error) {
	return i.repo.GetProfile(id)
}

func (i *ProfileInteractor) CreateProfile(data *profile.Create) (*profile.Created, error) {
	return i.repo.CreateProfile(data)
}

func (i *ProfileInteractor) UpdateProfile(id interface{}, data *profile.UpdateInfo) error {
	return i.repo.UpdateProfile(id, data)
}

func (i *ProfileInteractor) UpdateProfileAvatar(id, avatarPath interface{}) error {
	return i.repo.UpdateProfileAvatar(id, avatarPath)
}

func (i *ProfileInteractor) UpdateProfilePassword(id interface{}, data *profile.UpdatePassword) error {
	return i.repo.UpdateProfilePassword(id, data)
}

func (i *ProfileInteractor) GetProfileByEmail(email interface{}) (*profile.Profile, error) {
	return i.repo.GetProfileByEmail(email)
}

func (i *ProfileInteractor) GetProfileByNickname(nickname interface{}) (*profile.Profile, error) {
	return i.repo.GetProfileByNickname(nickname)
}

func (i *ProfileInteractor) GetProfileByEmailWithPassword(data *profile.Login) (*profile.Profile, error) {
	return i.repo.GetProfileByEmailWithPassword(data)
}

func (i *ProfileInteractor) GetProfiles(limit, offset int64) ([]profile.Info, error) {
	return i.repo.GetProfiles(limit, offset)
}

func (i *ProfileInteractor) GetProfileCount() (count int64, err error) {
	return i.repo.GetProfileCount()
}
