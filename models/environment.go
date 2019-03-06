package models

import (
	"api/database"
	"api/session"
)

type Env struct {
	Dbm *database.Manager
	Sm  *session.Manager
}
