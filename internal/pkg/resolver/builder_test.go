package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheme(t *testing.T) {
	r := &ResolverBuilder{}
	assert.Equal(t, r.Scheme(), Scheme, "Invalid scheme")
}
