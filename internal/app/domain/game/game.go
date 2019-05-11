package game

type State struct {
	Field       *Field            `json:"field,omitempty"`
	Players     map[string]Player `json:"players,omitempty"`
	ActiveItems map[uint64]Item   `json:"activeItems"`
	RoundNumber int               `json:"roundNumber,omitempty"`
	RoundTimer  *uint64           `json:"roundTimer,omitempty"`
}

// Empty checks state for emptiness
func (s *State) Empty() bool {
	if s.Field == nil && len(s.Players) == 0 && s.RoundNumber == 0 &&
		s.RoundTimer == nil && len(s.ActiveItems) == 0 {
		return true
	}

	return false
}

type Item struct {
	Type     string `json:"type"`
	PlayerId uint64 `json:"playerId"`
	Duration uint64 `json:"duration"`
}

type Field struct {
	Cells  []Cell `json:"cells,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type Cell struct {
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Type   string `json:"type"`
	HasBox bool   `json:"hasBox"`
}

//easyjson:json
type Player struct {
	Id        uint64            `json:"id"`
	X         int               `json:"x"`
	Y         int               `json:"y"`
	Items     map[string]uint64 `json:"items"`
	LoseRound *int              `json:"loseRound,omitempty"`
	Ready     bool              `json:"-"`
}

//easyjson:json
type Players []Player

//easyjson:json
type SetGameStartPayload struct {
	Field   *Field
	Players []Player
}

//easyjson:json
type InitPlayersPayload struct {
	PlayersId []uint64 `json:"playerIds"`
}

//easyjson:json
type InitPlayerMovePayload struct {
	PlayerId uint64 `json:"playerId"`
	Move     string `json:"move"`
}

//easyjson:json
type InitPlayerReadyPayload struct {
	PlayerId uint64 `json:"playerId"`
}

//easyjson:json
type InitItemUsePayload struct {
	PlayerId uint64 `json:"playerId"`
	ItemType string `json:"itemType"`
}

//easyjson:json
type SetItemPayload struct {
	PlayerId uint64 `json:"playerId"`
	ItemType string `json:"itemType"`
}
