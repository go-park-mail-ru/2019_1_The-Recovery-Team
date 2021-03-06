package game

import (
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/game"
)

type Transport struct {
	InnerReceiver func(action interface{})
	OuterReceiver func(action *game.Action)
}

// SendOut sends action out of the engine
func (t *Transport) SendOut(action *game.Action) {
	if t.OuterReceiver != nil {
		t.OuterReceiver(action)
	}
}

// SendInside sends action inside of the engine
func (t *Transport) SendInside(action interface{}) {
	if t.OuterReceiver != nil {
		t.InnerReceiver(action)
	}
}
