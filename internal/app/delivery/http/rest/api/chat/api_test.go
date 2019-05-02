package chat

import (
	"testing"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/memory/chat"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	sessionGrpc "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewRestApi(t *testing.T) {
	log, _ := zap.NewProduction()
	chatInteractor := usecase.NewChatInteractor(&chat.RepoMock{})
	sessionManager := sessionGrpc.NewClientMock()

	api := NewApi(chatInteractor, &sessionManager, log)
	assert.NotEmpty(t, api,
		"Returns empty api router")
}
