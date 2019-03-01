package handlers

const (
	QueryProfileById = `SELECT id, nickname, record, win, loss FROM profile WHERE id = $1`
)
