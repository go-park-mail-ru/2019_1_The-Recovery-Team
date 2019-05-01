package message

import (
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestNewChatRepo(t *testing.T) {
	conn := &pgx.ConnPool{}
	assert.NotEmpty(t, NewRepo(conn),
		"Doesn't create profile repository instance")
}
