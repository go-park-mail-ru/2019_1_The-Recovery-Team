package session

import (
	"time"

	"github.com/go-redis/redis"

	uuid "github.com/satori/go.uuid"
)

func NewSessionRepo(conn *redis.Client) *Repo {
	return &Repo{
		conn: conn,
	}
}

type Repo struct {
	conn *redis.Client
}

// Set sets value to redis
func (r *Repo) Set(profileID uint64, expires time.Duration) (string, error) {
	token := uuid.NewV4()

	if err := r.conn.Set(token.String(), profileID, expires).Err(); err != nil {
		return "", err
	}

	return token.String(), nil
}

// Delete deletes value from redis
func (r *Repo) Delete(token string) error {
	err := r.conn.Del(token).Err()
	return err
}

// Get gets value from redis
func (r *Repo) Get(token string) (uint64, error) {
	id, err := r.conn.Get(token).Uint64()
	if err != nil {
		return 0, err
	}

	return id, nil
}
