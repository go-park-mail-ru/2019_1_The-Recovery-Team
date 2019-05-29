package profile

import (
	"errors"
	"fmt"

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
	ORDER BY record DESC LIMIT $1 OFFSET $2`

	QueryProfileCount = `SELECT COUNT(*) 
	FROM profile`

	QueryProfileRatingById = `SELECT record 
	FROM profile 
	WHERE id = $1`

	QueryUpdateProfileRatingWinner = `UPDATE profile 
	SET record = record + $1, 
	win = win + 1
	WHERE id = $2`

	QueryUpdateProfileRatingLoser = `UPDATE profile 
	SET record = record - $1, 
	loss = loss + 1
	WHERE id = $2`

	QueryProfileIdOauth = `SELECT profile_id
	FROM token 
	WHERE user_id = $1`

	QueryCreateOauthToken = `INSERT INTO token (user_id, profile_id, token, oauth)
	VALUES ($1, $2, $3, $4)`

	QueryUpdateOauthToken = `UPDATE token 
	SET token = $1 
	WHERE user_id = $2`

	QueryCreateProfileOAuth = `INSERT INTO profile (nickname, avatar) 
	VALUES ($1, $2) 
	RETURNING id`

	QueryOauthByProfileId = `SELECT oauth, user_id 
	FROM token 
	WHERE profile_id = $1`

	QueryRatingPosition = `SELECT COUNT(*) FROM profile 
	WHERE record >= $1`

	NicknameAlreadyExists    = "NicknameAlreadyExists"
	EmailAlreadyExists       = "EmailAlreadyExists"
	IncorrectProfilePassword = "IncorrectProfilePassword"

	ProfileEmailKey    = "profile_email_key"
	ProfileNicknameKey = "profile_nickname_key"

	RatingStep            = 60
	DefaultRatingIncrease = 25
	MinRatingIncrease     = 1
)

// NewRepo creates new instance of profile repository
func NewRepo(conn *pgx.ConnPool) *Repo {
	return &Repo{
		conn: postgresql.NewConnPool(conn),
	}
}

type Repo struct {
	conn postgresql.Conn
}

// Get gets profile data by id
func (r *Repo) Get(id interface{}) (*profile.Profile, error) {
	profile := &profile.Profile{}
	var email *string
	if err := r.conn.QueryRow(QueryProfileById, id).Scan(&profile.ID, &profile.Nickname, &email,
		&profile.Avatar, &profile.Record, &profile.Win, &profile.Loss); err != nil {
		return nil, err
	}

	if email != nil {
		profile.Email = *email
	}

	if err := r.conn.QueryRow(QueryRatingPosition, profile.Record).Scan(&profile.Position); err != nil {
		return nil, err
	}

	err := r.conn.QueryRow(QueryOauthByProfileId, id).Scan(&profile.Oauth, &profile.OauthId)
	if err != nil && err != pgx.ErrNoRows {
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
	var emailReceive *string
	err := r.conn.QueryRow(QueryProfileByEmail, email).Scan(&received.ID, &emailReceive, &received.Nickname)
	if emailReceive != nil {
		received.Email = *emailReceive
	}
	return received, err
}

// GetByNickname gets profile by nickname
func (r *Repo) GetByNickname(nickname interface{}) (*profile.Profile, error) {
	received := &profile.Profile{}
	var email *string
	err := r.conn.QueryRow(QueryProfileByNickname, nickname).Scan(&received.ID, &email, &received.Nickname)
	if email != nil {
		received.Email = *email
	}
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
		prof := profile.Info{}
		if err = rows.Scan(&prof.ID, &prof.Nickname, &prof.Avatar, &prof.Record, &prof.Win, &prof.Loss); err != nil {
			return nil, err
		}
		profiles = append(profiles, prof)
	}
	return profiles, nil
}

// Count gets number of profiles
func (r *Repo) Count() (count int64, err error) {
	err = r.conn.QueryRow(QueryProfileCount).Scan(&count)
	return
}

// UpdateRating updates players rating
func (r *Repo) UpdateRating(winner, loser uint64) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Getting players rating
	var winnerRating, loserRating int64
	if err := tx.QueryRow(QueryProfileRatingById, winner).Scan(&winnerRating); err != nil {
		return err
	}
	if err := tx.QueryRow(QueryProfileRatingById, loser).Scan(&loserRating); err != nil {
		return err
	}

	// Calculating rating increase
	increase := DefaultRatingIncrease + (loserRating-winnerRating)/RatingStep
	if increase <= 0 {
		increase = MinRatingIncrease
	}

	if _, err := tx.Exec(QueryUpdateProfileRatingLoser, increase, loser); err != nil {
		return err
	}
	if _, err := tx.Exec(QueryUpdateProfileRatingWinner, increase, winner); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (r *Repo) PutProfileOauth(id string, token string) (*profile.ID, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var profileId uint64
	if err := tx.QueryRow(QueryProfileIdOauth, id).Scan(&profileId); err != nil {
		return nil, errors.New("ProfileDoesNotExist")
	}

	if _, err := tx.Exec(QueryUpdateOauthToken, token, id); err != nil {
		return nil, err
	}

	tx.Commit()
	return &profile.ID{Id: profileId}, nil
}

func (r *Repo) CreateProfileOauth(create *profile.CreateOauth) (*profile.ID, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var profileId uint64
	if err := tx.QueryRow(QueryCreateProfileOAuth, fmt.Sprintf("%s_%s", create.Oauth, create.UserId), create.Avatar.Path).Scan(&profileId); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(QueryCreateOauthToken, create.UserId, profileId, create.Token, create.Oauth); err != nil {
		return nil, err
	}

	tx.Commit()
	return &profile.ID{Id: profileId}, nil
}
