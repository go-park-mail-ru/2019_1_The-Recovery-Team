package game

import (
	"sync"
	"time"
)

type Action struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type State struct {
	Field       *Field            `json:"field,omitempty"`
	Players     map[string]Player `json:"players,omitempty"`
	ActiveItems *sync.Map         `json:"active_items,omitempty"`
	RoundNumber int               `json:"round_number,omitempty"`
	RoundTimer  uint64            `json:"round_timer,omitempty"`
}

type Item struct {
	Type     string        `json:"type"`
	PlayerId uint64        `json:"player_id"`
	Duration time.Duration `json:"duration"`
}

type Field struct {
	Cells  []Cell `json:"cells"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Cell struct {
	Row    int    `json:"row"`
	Col    int    `json:"col"`
	Type   string `json:"type"`
	HasBox bool   `json:"has_box"`
}

//easyjson:json
type Player struct {
	Id        uint64            `json:"id"`
	X         int               `json:"x"`
	Y         int               `json:"y"`
	Items     map[string]uint64 `json:"items"`
	LoseRound *int              `json:"lose_round,omitempty"`
}

//easyjson:json
type Players []Player

//easyjson:json
type GameStartPayload struct {
	Field   *Field
	Players []Player
}

//easyjson:json
type InitPlayersPayload struct {
	PlayersId []uint64 `json:"players"`
}

type InitPlayerMovePayload struct {
	PlayerId uint64 `json:"player_id"`
	Move     string `json:"move"`
}
