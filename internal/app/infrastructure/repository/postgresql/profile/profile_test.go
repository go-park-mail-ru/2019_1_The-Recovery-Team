package profile

import (
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestNewProfileRepo(t *testing.T) {
	conn := &pgx.Conn{}
	assert.NotEmpty(t, NewRepo(conn),
		"Doesn't create profile repository instance")
}
