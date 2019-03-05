package models

import (
	"api/database"
)

type Env struct {
	Dbm *database.Manager
}
