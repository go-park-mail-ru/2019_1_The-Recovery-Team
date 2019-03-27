package repository

import "time"

type SessionRepo interface {
	Get(token string) (uint64, error)
	Set(profileID uint64, expires time.Duration) (string, error)
	Delete(token string) error
}
