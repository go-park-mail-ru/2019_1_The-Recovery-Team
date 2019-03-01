package models

import (
	"api/database"
)

// Env containts environment variables
type Env struct {
	Dbm *database.Manager
}
