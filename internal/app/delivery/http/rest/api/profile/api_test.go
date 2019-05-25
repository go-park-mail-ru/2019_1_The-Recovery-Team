package profile

import (
	"testing"

	sessionGrpc "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"

	profileGrpc "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewRestApi(t *testing.T) {
	log, _ := zap.NewProduction()
	profileManager := profileGrpc.NewClientMock()
	sessionManager := sessionGrpc.NewClientMock()

	api := NewApi(&profileManager, &sessionManager, log, "", "")
	assert.NotEmpty(t, api,
		"Returns empty api router")
}
