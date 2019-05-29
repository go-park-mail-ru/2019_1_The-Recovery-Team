package session

import (
	"testing"
	"time"

	"github.com/go-redis/redis"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func newMockRedis() *redis.Client {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}

	return client
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
