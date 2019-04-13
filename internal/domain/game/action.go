package game

const (
	SetGameStart        = "SET_GAME_START"
	InitPlayers         = "INIT_PLAYERS"
	InitPlayerMove      = "INIT_PLAYER_MOVE"
	SetRoundStart       = "SET_ROUND_START"
	SetRoundTime        = "SET_ROUND_TIME"
	SetState            = "SET_STATE"
	SetStateDiff        = "SET_STATE_DIFF"
	SetRoundStop        = "SET_ROUND_STOP"
	InitPlayerReady     = "INIT_PLAYER_READY"
	SetFieldRound       = "SET_FIELD_ROUND"
	SetGameOver         = "SET_GAME_OVER"
	SetUserDisconnected = "SET_USER_DISCONNECTED"
	SetOpponentLeave    = "SET_OPPONENT_LEAVE"
	InitEngineStop      = "INIT_ENGINE_STOP"
	SetEngineStop       = "SET_ENGINE_STOP"
	SetOpponentSearch   = "SET_OPPONENT_SEARCH"
	SetAlreadyPlaying   = "SET_ALREADY_PLAYING"
	SetOpponentNotFound = "SET_OPPONENT_NOT_FOUND"
	SetOpponent         = "SET_OPPONENT"
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
