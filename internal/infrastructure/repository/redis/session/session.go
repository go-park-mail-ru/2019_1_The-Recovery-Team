package session

import (
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

func NewSessionRepo(conn *redis.Conn) *Repo {
	return &Repo{
		conn: conn,
	}
}

type Repo struct {
	conn *redis.Conn
}

// Set sets value to redis
func (r *Repo) Set(profileID uint64, expires time.Duration) (string, error) {
	token := uuid.NewV4()

	_, err := (*r.conn).Do("SETEX", token.String(), int(expires), profileID)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

// Delete deletes value from redis
func (r *Repo) Delete(token string) error {
	_, err := (*r.conn).Do("DEL", token)
	return err
}

// Get gets value from redis
func (r *Repo) Get(token string) (uint64, error) {
	id, err := redis.Uint64((*r.conn).Do("GET", token))
	if err != nil {
		return 0, err
	}

	return id, nil
}
