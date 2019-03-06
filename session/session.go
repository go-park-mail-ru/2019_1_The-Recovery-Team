package session

import (
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

// Manager is session manager
type Manager struct {
	sm redis.Conn
}

// InitSessionManager initialize session manager
func InitSessionManager(username, password, host string) (*Manager, error) {
	var err error
	manager := &Manager{}
	manager.sm, err = redis.DialURL("redis://" + username + ":" + password + "@" + host)
	if err != nil {
		return manager, err
	}
	return manager, nil
}

// SM returns connection to redis
func (manager *Manager) SM() *redis.Conn {
	return &manager.sm
}

// Set sets value to redis
func (manager *Manager) Set(profileID uint64, expires time.Duration) (string, error) {
	token := uuid.NewV4()

	_, err := manager.sm.Do("SETEX", token.String(), int(expires), profileID)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

// Delete deletes value from redis
func (manager *Manager) Delete(token string) error {
	_, err := manager.sm.Do("DEL", token)
	return err
}

// Get gets value from redis
func (manager *Manager) Get(token string) (uint64, error) {
	id, err := redis.Uint64(manager.sm.Do("GET", token))
	if err != nil {
		return 0, err
	}

	return id, nil
}
