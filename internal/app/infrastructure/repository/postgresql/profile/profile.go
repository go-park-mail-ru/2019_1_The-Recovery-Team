package profile

import (
	"errors"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql"

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

	QueryProfileCount = `SELECT reltuples::bigint AS number
	FROM   pg_class
	WHERE  oid = 'public.profile'::regclass`

	NicknameAlreadyExists    = "NicknameAlreadyExists"
	EmailAlreadyExists       = "EmailAlreadyExists"
	IncorrectProfilePassword = "IncorrectProfilePassword"

	ProfileEmailKey    = "profile_email_key"
	ProfileNicknameKey = "profile_nickname_key"
)

// NewRepo creates new instance of profile repository
func NewRepo(conn *pgx.Conn) *Repo {
	return &Repo{
		conn: conn,
	}
}

type Repo struct {
	conn *pgx.Conn
}

// Get gets profile data by id
func (r *Repo) Get(id interface{}) (*profile.Profile, error) {
	profile := &profile.Profile{}
	if err := r.conn.QueryRow(QueryProfileById, id).Scan(&profile.ID, &profile.Nickname, &profile.Email,
		&profile.Avatar, &profile.Record, &profile.Win, &profile.Loss); err != nil {
		return nil, err
	}

	return profile, nil
}

// Create creates new profile
func (r *Repo) Create(data *profile.Create) (*profile.Created, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	created := &profile.Created{}
	if err = tx.QueryRow(QueryCreateProfile, data.Email, data.Nickname, data.Password).
		Scan(&created.ID, &created.Email, &created.Nickname); err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			switch pgErr.ConstraintName {
			case ProfileEmailKey:
				{
					return nil, errors.New(EmailAlreadyExists)
				}
			case ProfileNicknameKey:
				{
					return nil, errors.New(NicknameAlreadyExists)
				}
			}
		}
		return nil, err
	}

	tx.Commit()
	return created, nil
}

// Update updates profile
func (r *Repo) Update(id interface{}, data *profile.UpdateInfo) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(QueryUpdateProfile, data.Email, data.Nickname, id); err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			switch pgErr.ConstraintName {
			case ProfileEmailKey:
				{
					return errors.New(EmailAlreadyExists)
				}
			case ProfileNicknameKey:
				{
					return errors.New(NicknameAlreadyExists)
				}
			}
		}
		return err
	}

	tx.Commit()
	return nil
}

// UpdateAvatar updates profile avatar
func (r *Repo) UpdateAvatar(id, avatarPath interface{}) error {
	_, err := r.conn.Exec(QueryUpdateProfileAvatar, avatarPath, id)
	return err
}

// UpdatePassword updates profile password
func (r *Repo) UpdatePassword(id interface{}, data *profile.UpdatePassword) error {
	var password string
	if err := r.conn.QueryRow(QueryProfileByIdWithPassword, id).
		Scan(&password); err != nil {
		return err
	}

	// Check current password correctness
	if matches, err := postgresql.VerifyPassword(data.PasswordOld, password); !matches || err != nil {
		return errors.New(IncorrectProfilePassword)
	}

	_, err := r.conn.Exec(QueryUpdateProfilePassword, data.Password, id)
	return err
}

// GetByEmail gets profile by email
func (r *Repo) GetByEmail(email interface{}) (*profile.Profile, error) {
	received := &profile.Profile{}
	err := r.conn.QueryRow(QueryProfileByEmail, email).Scan(&received.ID, &received.Email, &received.Nickname)
	return received, err
}

// GetByNickname gets profile by nickname
func (r *Repo) GetByNickname(nickname interface{}) (*profile.Profile, error) {
	received := &profile.Profile{}
	err := r.conn.QueryRow(QueryProfileByNickname, nickname).Scan(&received.ID, &received.Email, &received.Nickname)
	return received, err
}

// GetByEmailAndPassword gets profile by email and password(login)
func (r *Repo) GetByEmailAndPassword(data *profile.Login) (*profile.Profile, error) {
	received := &profile.Profile{}
	if err := r.conn.QueryRow(QueryProfileByEmailWithPassword, data.Email).
		Scan(&received.ID, &received.Email, &received.Nickname, &received.Password, &received.Avatar, &received.Record, &received.Win, &received.Loss); err != nil {
		return nil, err
	}

	if matches, err := postgresql.VerifyPassword(data.Password, received.Password); !matches || err != nil {
		return nil, pgx.ErrNoRows
	}

	return received, nil
}

// List gets profile list
func (r *Repo) List(limit, offset int64) ([]profile.Info, error) {
	profiles := make([]profile.Info, 0, 10)
	rows, err := r.conn.Query(QueryProfilesWithLimitAndOffset, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		profile := profile.Info{}
		if err = rows.Scan(&profile.ID, &profile.Nickname, &profile.Avatar, &profile.Record, &profile.Win, &profile.Loss); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// Count gets number of profiles
func (r *Repo) Count() (count int64, err error) {
	err = r.conn.QueryRow(QueryProfileCount).Scan(&count)
	return
}
