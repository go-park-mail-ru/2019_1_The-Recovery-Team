package session

import (
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

type Manager struct {
	sm redis.Conn
}

func InitSessionManager(username, password, host string) (*Manager, error) {
	var err error
	manager := &Manager{}
	manager.sm, err = redis.DialURL("redis://" + username + ":" + password + "@" + host)
	if err != nil {
		return manager, err
	}
	return manager, nil
}

func (manager *Manager) SM() *redis.Conn {
	return &manager.sm
}

func (manager *Manager) Set(profileID uint64, expires time.Duration) (string, error) {
	token := uuid.NewV4()

	_, err := manager.sm.Do("SETEX", token.String(), int(expires), profileID)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

func (manager *Manager) Delete(token string) error {
	_, err := manager.sm.Do("DEL", token)
	return err
}

func (manager *Manager) Get(token string) (uint64, error) {
	id, err := redis.Uint64(manager.sm.Do("GET", token))
	if err != nil {
		return 0, err
	}

	return id, nil
}
