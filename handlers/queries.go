package handlers

const (
	QueryProfiles                   = `SELECT id, nickname, avatar, record, win, loss FROM profile ORDER BY record`
	QueryProfilesWithLimit          = `SELECT id, nickname, avatar, record, win, loss FROM profile ORDER BY record LIMIT $1`
	QueryProfilesWithOffset         = `SELECT id, nickname, avatar, record, win, loss FROM profile ORDER BY record OFFSET $1`
	QueryProfilesWithLimitAndOffset = `SELECT id, nickname, avatar, record, win, loss FROM profile ORDER BY record LIMIT $1 OFFSET $2`
	QueryProfileById                = `SELECT id, nickname, email, avatar, record, win, loss FROM profile WHERE id = $1`
	QueryProfileByEmail             = `SELECT email FROM profile WHERE LOWER(email) = LOWER($1)`
	QueryProfileByNickname          = `SELECT nickname FROM profile WHERE LOWER(nickname) = LOWER($1)`
	QueryInsertProfile              = `INSERT INTO profile (email, nickname, password) VALUES ($1, $2, $3) RETURNING id, email, nickname`
	QueryProfileUnsafe              = `SELECT id, email, nickname, record, win, loss FROM profile WHERE id = $1`
	QueryUpdateProfileAvatar        = `UPDATE profile SET avatar = $1 WHERE id = $2`
	QueryProfileByEmailWithPassword = `SELECT id, email, nickname, password, avatar FROM profile WHERE LOWER(email) = LOWER($1)`
	QueryCountProfilesNumber        = `SELECT COUNT(*) FROM profile`
)
