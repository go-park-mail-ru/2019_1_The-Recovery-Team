package game

type Transport struct {
	InnerReceiver func(action string)
	OuterReceiver func(actionType, action string)
}

func (t *Transport) SendOut(actionType string, action string) {
	if t.OuterReceiver != nil {
		t.OuterReceiver(actionType, action)
	}
}

func (t *Transport) SendInside(action string) {
	if t.OuterReceiver != nil {
		t.InnerReceiver(action)
	}
}
