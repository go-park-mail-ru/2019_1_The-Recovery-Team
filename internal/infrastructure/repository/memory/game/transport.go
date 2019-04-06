package game

import "sadislands/internal/domain/game"

type Transport struct {
	InnerReceiver func(action interface{})
	OuterReceiver func(action *game.Action)
}

func (t *Transport) SendOut(action *game.Action) {
	if t.OuterReceiver != nil {
		t.OuterReceiver(action)
	}
}

func (t *Transport) SendInside(action interface{}) {
	if t.OuterReceiver != nil {
		t.InnerReceiver(action)
	}
}
