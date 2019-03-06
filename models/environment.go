package models

import (
	"api/database"
	"api/session"
)

// Env environmet model
type Env struct {
	Dbm *database.Manager
	Sm  *session.Manager
}
