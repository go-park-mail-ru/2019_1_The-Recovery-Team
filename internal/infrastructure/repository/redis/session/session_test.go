package session

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func newMockRedis() *redis.Conn {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	conn, err := redis.Dial("tcp", s.Addr())
	return &conn
}

func TestNewSessionRepo(t *testing.T) {
	assert.NotEmpty(t, NewSessionRepo(newMockRedis()),
		"Doesn't create session repository instance")
}

func TestSet(t *testing.T) {
	repo := NewSessionRepo(newMockRedis())
	_, err := repo.Set(DefaultProfileId, time.Second)
	assert.Empty(t, err, "Doesn't correctly set session")
}

func TestGet(t *testing.T) {
	repo := NewSessionRepo(newMockRedis())
	token, _ := repo.Set(DefaultProfileId, time.Second)
	profileId, err := repo.Get(token)
	assert.Empty(t, err, "Doesn't correctly get session")
	assert.Equal(t, DefaultProfileId, profileId)
}

func TestDelete(t *testing.T) {
	repo := NewSessionRepo(newMockRedis())
	token, _ := repo.Set(DefaultProfileId, time.Second)
	err := repo.Delete(token)
	assert.Empty(t, err, "Doesn't correctly delete session")
}
