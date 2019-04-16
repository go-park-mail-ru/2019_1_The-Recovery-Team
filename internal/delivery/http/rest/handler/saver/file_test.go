package saver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashFileName(t *testing.T) {
	hash, err := HashFileName("filename", 1)
	assert.Empty(t, err,
		"Returns error on correct data")
	assert.NotEmpty(t, hash,
		"Doesn't create hash string")
}
