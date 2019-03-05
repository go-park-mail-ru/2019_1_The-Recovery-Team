package handlers

const (
	QueryProfiles                   = `SELECT nickname, avatar, record, win, loss FROM profile`
	QueryProfilesWithLimit          = `SELECT nickname, avatar, record, win, loss FROM profile LIMIT $1`
	QueryProfilesWithOffset         = `SELECT nickname, avatar, record, win, loss FROM profile OFFSET $1`
	QueryProfilesWithLimitAndOffset = `SELECT nickname, avatar, record, win, loss FROM profile LIMIT $1 OFFSET $2 ORDER BY score`
	QueryProfileById                = `SELECT nickname, avatar, record, win, loss FROM profile WHERE id = $1`
	QueryProfileByEmail             = `SELECT email FROM profile WHERE LOWER(email) = LOWER($1)`
	QueryProfileByNickname          = `SELECT nickname FROM profile WHERE LOWER(nickname) = LOWER($1)`
	QueryInsertProfile              = `INSERT INTO profile (email, nickname, password) VALUES ($1, $2, $3) RETURNING id, email, nickname`
	QueryProfileUnsafe              = `SELECT id, email, nickname, record, win, loss FROM profile WHERE id = $1`
	QueryUpdateProfileAvatar        = `UPDATE profile SET avatar = $1 WHERE id = $2`
)
