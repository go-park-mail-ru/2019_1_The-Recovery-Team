package chat

const (
	SetMessage         = "SET_CHAT_MESSAGE"
	InitMessage        = "INIT_CHAT_MESSAGE"
	InitPing           = "INIT_PING"
	SetPong            = "SET_PONG"
	SetSession         = "SET_CHAT_SESSION"
	InitGlobalMessages = "INIT_CHAT_GLOBAL_MESSAGES"
	SetGlobalMessages  = "SET_CHAT_GLOBAL_MESSAGES"
	InitUpdateMessage  = "INIT_CHAT_MESSAGE_UPDATE"
	SetUpdateMessage   = "SET_CHAT_MESSAGE_UPDATE"
	InitPrinting       = "INIT_CHAT_PRINTING"
	SetPrinting        = "SET_CHAT_PRINTING"
	InitDeleteMessage  = "INIT_CHAT_MESSAGE_DELETE"
	SetDeleteMessage   = "INIT_CHAT_MESSAGE_DELETE"
)

//easyjson:json
type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

//easyjson:json
type ActionRaw struct {
	Type    string `json:"type"`
	Payload string `json:"payload,omitempty"`
}

//{"type": "INIT_CHAT_MESSAGE_UPDATE", "payload": "{\"messageId\":24, \"data\":{\"text\":\"Fedya Pidor\"}}"}
