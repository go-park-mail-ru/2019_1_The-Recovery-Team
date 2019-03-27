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

func (s *ProfileInteractor) GetProfile(id interface{}) (*profile.Profile, error) {
	return s.repo.GetProfile(id)
}

func (s *ProfileInteractor) CreateProfile(data *profile.Create) (*profile.Created, error) {
	return s.repo.CreateProfile(data)
}

func (s *ProfileInteractor) UpdateProfile(id interface{}, data *profile.UpdateInfo) error {
	return s.repo.UpdateProfile(id, data)
}

func (s *ProfileInteractor) UpdateProfileAvatar(id, avatarPath interface{}) error {
	return s.repo.UpdateProfileAvatar(id, avatarPath)
}

func (s *ProfileInteractor) UpdateProfilePassword(id interface{}, data *profile.UpdatePassword) error {
	return s.repo.UpdateProfilePassword(id, data)
}

func (s *ProfileInteractor) GetProfileByEmail(email interface{}) (*profile.Profile, error) {
	return s.repo.GetProfileByEmail(email)
}

func (s *ProfileInteractor) GetProfileByNickname(nickname interface{}) (*profile.Profile, error) {
	return s.repo.GetProfileByNickname(nickname)
}

func (s *ProfileInteractor) GetProfileByEmailWithPassword(data *profile.Login) (*profile.Profile, error) {
	return s.repo.GetProfileByEmailWithPassword(data)
}

func (s *ProfileInteractor) GetProfiles(limit, offset int64) ([]profile.Info, error) {
	return s.repo.GetProfiles(limit, offset)
}

func (s *ProfileInteractor) GetProfileCount() (count int64, err error) {
	return s.repo.GetProfileCount()
}
