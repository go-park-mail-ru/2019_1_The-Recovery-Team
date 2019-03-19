package environment

import (
	"api/database"
	"api/session"
	"go.uber.org/zap"
)

// Env environmet model
type Env struct {
	Dbm *database.Manager
	Sm  *session.Manager
	Lm  *zap.Logger
}
