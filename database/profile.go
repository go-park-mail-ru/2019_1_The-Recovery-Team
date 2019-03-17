package database

import (
	"api/models"
	"github.com/jackc/pgx"
)

const (
	QueryProfileById = `SELECT id, nickname, email, avatar, record, win, loss 
    FROM profile 
    WHERE id = $1`

	QueryCreateProfile = `INSERT INTO profile (email, nickname, password) 
	VALUES ($1, $2, $3) 
	RETURNING id, email, nickname`

	QueryUpdateProfile = `UPDATE profile
	SET email = (CASE WHEN $1 = '' THEN email ELSE $1 END),
	nickname = (CASE WHEN $2 = '' THEN nickname ELSE $2 END)
	WHERE id = $3`

	QueryUpdateProfileAvatar = `UPDATE profile
	SET avatar = $1
	WHERE id = $2`

	QueryUpdateProfilePassword = `UPDATE profile
	SET password = $1
	WHERE id = $2`

	QueryProfileByEmail = `SELECT id, email, nickname
	FROM profile
	WHERE email = $1`

	QueryProfileByNickname = `SELECT id, email, nickname
	FROM profile
	WHERE nickname = $1`

	QueryProfileByIdWithPassword = `SELECT password
	FROM profile
	WHERE id = $1`

	QueryProfileByEmailWithPassword = `SELECT id, email, nickname, password, avatar, record, win, loss  
	FROM profile 
	WHERE email = $1`

	QueryProfilesWithLimitAndOffset = `SELECT id, nickname, avatar, record, win, loss 
	FROM profile 
	ORDER BY record LIMIT $1 OFFSET $2`

	QueryProfilesNumber = `SELECT reltuples::bigint AS number
	FROM   pg_class
	WHERE  oid = 'public.profile'::regclass`
)

func (dbm *Manager) GetProfile(id interface{}) (*models.Profile, error) {
	profile := &models.Profile{}
	err := dbm.conn.QueryRow(QueryProfileById, id).Scan(&profile.ID, &profile.Nickname, &profile.Email,
		&profile.Avatar, &profile.Record, &profile.Win, &profile.Loss)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (dbm *Manager) CreateProfile(data *models.ProfileCreate) (*models.ProfileCreated, error) {
	tx, err := dbm.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	created := &models.ProfileCreated{}
	err = tx.QueryRow(QueryCreateProfile, data.Email, data.Nickname, data.Password).Scan(&created.ID, &created.Email, &created.Nickname)
	if err != nil {
		switch {
		case err.(pgx.PgError).ConstraintName == "profile_email_key":
			{
				return nil, ErrEmailAlreadyExists
			}
		case err.(pgx.PgError).ConstraintName == "profile_nickname_key":
			{
				return nil, ErrNicknameAlreadyExists
			}
		default:
			{
				return nil, err
			}
		}
	}

	tx.Commit()
	return created, nil
}

func (dbm *Manager) UpdateProfile(id interface{}, data *models.ProfileUpdate) error {
	tx, err := dbm.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(QueryUpdateProfile, data.Email, data.Nickname, id)
	if err != nil {
		switch {
		case err.(pgx.PgError).ConstraintName == "profile_email_key":
			{
				return ErrEmailAlreadyExists
			}
		case err.(pgx.PgError).ConstraintName == "profile_nickname_key":
			{
				return ErrNicknameAlreadyExists
			}
		default:
			{
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (dbm *Manager) UpdateProfileAvatar(id, avatarPath interface{}) error {
	_, err := dbm.conn.Exec(QueryUpdateProfileAvatar, avatarPath, id)
	return err
}

func (dbm *Manager) UpdateProfilePassword(id interface{}, data *models.ProfileUpdatePassword) error {
	var password string
	err := dbm.conn.QueryRow(QueryProfileByIdWithPassword, id).Scan(&password)
	if err != nil {
		return err
	}
	if matches, err := VerifyPassword(data.PasswordOld, password); !matches || err != nil {
		return ErrIncorrectPassword
	}

	_, err = dbm.conn.Exec(QueryUpdateProfilePassword, data.Password, id)
	return err
}

func (dbm *Manager) GetProfileByEmail(email interface{}) (*models.Profile, error) {
	profile := &models.Profile{}
	err := dbm.conn.QueryRow(QueryProfileByEmail, email).Scan(&profile.ID, &profile.Email, &profile.Nickname)
	return profile, err
}

func (dbm *Manager) GetProfileByNickname(nickname interface{}) (*models.Profile, error) {
	profile := &models.Profile{}
	err := dbm.conn.QueryRow(QueryProfileByNickname, nickname).Scan(&profile.ID, &profile.Email, &profile.Nickname)
	return profile, err
}

func (dbm *Manager) GetProfileByEmailWithPassword(data *models.ProfileLogin) (*models.Profile, error) {
	profile := &models.Profile{}
	err := dbm.conn.QueryRow(QueryProfileByEmailWithPassword, data.Email).Scan(&profile.ID, &profile.Email, &profile.Nickname, &profile.Password, &profile.Avatar, &profile.Record, &profile.Win, &profile.Loss)
	if err != nil {
		return nil, err
	}

	if matches, err := VerifyPassword(data.Password, profile.Password); !matches || err != nil {
		return nil, pgx.ErrNoRows
	}

	return profile, nil
}

func (dbm *Manager) GetProfiles(limit, offset int64) ([]models.ProfileInfo, error) {
	profiles := make([]models.ProfileInfo, 0, 10)
	rows, err := dbm.conn.Query(QueryProfilesWithLimitAndOffset, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		profile := models.ProfileInfo{}
		err = rows.Scan(&profile.ID, &profile.Nickname, &profile.Avatar, &profile.Record, &profile.Win, &profile.Loss)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (dbm *Manager) GetProfilesNumber() (number int64, err error) {
	err = dbm.conn.QueryRow(QueryProfilesNumber).Scan(&number)
	return
}
