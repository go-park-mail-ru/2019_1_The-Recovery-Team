package game

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"math/rand"
	"sadislands/internal/domain/game"
	"strconv"
	"sync"
	"time"
)

const (
	Sand                    = "SAND"
	Water                   = "WATER"
	Swamp                   = "SWAMP"
	WaterStartNumber        = 10
	SwampStartNumber        = 15
	FieldWidth              = 10
	FieldHeight             = 10
	CellSize                = 4
	Up                      = "UP"
	Down                    = "DOWN"
	Left                    = "LEFT"
	Right                   = "RIGHT"
	RoundDuration    uint64 = 5
	TickerDuration          = 1000 / 60
)

type Engine struct {
	Transport *Transport
	State     *game.State
	StateDiff *game.State

	UpdateM *sync.Mutex

	ProcessActions []*game.Action
	ProcessM       *sync.Mutex

	RoundRunning *atomic.Value
	GameOver     *atomic.Value

	ReceivedActions chan *game.Action

	Ticker *time.Ticker
	Timer  *time.Timer
}

func initState() *game.State {
	return &game.State{
		Field:       initField(),
		Players:     make(map[string]game.Player),
		ActiveItems: sync.Map{},
		RoundNumber: 0,
	}
}

func initField() *game.Field {
	field := &game.Field{
		Cells:  make([]game.Cell, 0, FieldWidth*FieldHeight),
		Width:  FieldWidth,
		Height: FieldHeight,
	}

	types := make([]string, 0, FieldWidth*FieldHeight)
	for i := 0; i < SwampStartNumber; i++ {
		types = append(types, Swamp)
	}
	for i := 0; i < WaterStartNumber; i++ {
		types = append(types, Water)
	}
	for i := 0; i < FieldWidth*FieldHeight-SwampStartNumber-WaterStartNumber; i++ {
		types = append(types, Sand)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(types), func(i, j int) { types[i], types[j] = types[j], types[i] })
	types[0] = Sand

	for i := 0; i < FieldHeight; i++ {
		for j := 0; j < FieldWidth; j++ {
			cell := game.Cell{
				Row:    i,
				Col:    j,
				Type:   types[i+j*FieldWidth],
				HasBox: false,
			}

			field.Cells = append(field.Cells, cell)
		}
	}

	return field
}

// Game action handlers
func (e *Engine) setGameStart() {
	e.State = initState()
	//e.StateDiff = e.State
}

// Round action handlers
func (e *Engine) setRoundStart() {
	e.State.RoundNumber += 1
	time.AfterFunc(time.Duration(RoundDuration)*time.Second, e.stopRound)

	e.State.RoundTimer = new(uint64)
	e.StateDiff.RoundTimer = new(uint64)
	atomic.StoreUint64(e.State.RoundTimer, RoundDuration)

	go e.controlRemainingRoundTime()

	e.StateDiff.RoundNumber = e.State.RoundNumber

	atomic.StoreUint64(e.StateDiff.RoundTimer, atomic.LoadUint64(e.State.RoundTimer))
}

func (e *Engine) setRoundTime() {
	e.StateDiff.RoundTimer = new(uint64)
	atomic.StoreUint64(e.State.RoundTimer, atomic.LoadUint64(e.State.RoundTimer)-1)
	atomic.StoreUint64(e.StateDiff.RoundTimer, atomic.LoadUint64(e.State.RoundTimer))
}

func (e *Engine) setRoundStop() {
	for idStr, player := range e.State.Players {
		player.Ready = false
		e.State.Players[idStr] = player
	}

	e.StateDiff.RoundTimer = new(uint64)
	atomic.StoreUint64(e.State.RoundTimer, 0)
	atomic.StoreUint64(e.StateDiff.RoundTimer, atomic.LoadUint64(e.State.RoundTimer))
}

// Players action handlers
func (e *Engine) initPlayers(action *game.Action) {
	if len(e.State.Players) != 0 {
		return
	}

	payload := action.Payload.(*game.InitPlayersPayload)

	for _, id := range payload.PlayersId {
		idStr := strconv.FormatUint(id, 10)

		e.State.Players[idStr] = game.Player{
			Id:    id,
			Items: make(map[string]uint64),
		}
	}
	//e.StateDiff.Players = e.State.Players
}

func (e *Engine) initPlayerReady(action *game.Action) {
	payload := action.Payload.(*game.InitPlayerReadyPayload)

	idStr := strconv.FormatUint(payload.PlayerId, 10)
	player, exists := e.State.Players[idStr]

	if !exists {
		return
	}

	if player.Ready {
		return
	}

	player.Ready = true
	e.State.Players[idStr] = player

	for _, player := range e.State.Players {
		if !player.Ready {
			return
		}
	}

	e.ReceivedActions <- &game.Action{
		Type: game.SetRoundStart,
	}
}

func (e *Engine) movePlayer(action *game.Action) {
	payload := action.Payload.(*game.InitPlayerMovePayload)

	idStr := strconv.FormatUint(payload.PlayerId, 10)
	player, exists := e.State.Players[idStr]

	if !exists || player.LoseRound != nil {
		return
	}

	switch payload.Move {
	case Right:
		{
			if player.X == FieldWidth*CellSize-1 {
				return
			}

			player.X += 1
		}
	case Up:
		{
			if player.Y == 0 {
				return
			}

			player.Y -= 1
		}
	case Left:
		{
			if player.X == 0 {
				return
			}

			player.X -= 1
		}
	case Down:
		{
			if player.Y == FieldHeight*CellSize-1 {
				return
			}

			player.Y += 1
		}
	default:
		{
			return
		}
	}

	e.State.Players[idStr] = player

	if e.isPlayerDead(idStr) {
		player.LoseRound = &e.State.RoundNumber
		e.GameOver.Store(true)
	}

	e.StateDiff.Players[idStr] = player
}

func (e *Engine) isPlayerDead(id string) bool {
	player := e.State.Players[id]
	if e.State.Field.Cells[(player.X/CellSize)+(player.Y/CellSize)*FieldWidth].Type == Water {
		return true
	}

	return false
}

func (e *Engine) controlRemainingRoundTime() {
	action := &game.Action{
		Type: game.SetRoundTime,
	}

	t := time.NewTicker(1 * time.Second)
	for {
		<-t.C

		if atomic.LoadUint64(e.State.RoundTimer) == 2 {
			e.ReceivedActions <- action
			t.Stop()
			return
		}

		e.ReceivedActions <- action
	}
}

func (e *Engine) stopRound() {
	e.ReceivedActions <- &game.Action{
		Type: game.SetRoundStop,
	}
}

func initStateDiff() *game.State {
	return &game.State{
		Players: make(map[string]game.Player),
	}
}

func (e *Engine) copyState() *game.State {
	state := &game.State{
		Field: &game.Field{},
		Players: e.State.Players,
		ActiveItems: e.State.ActiveItems,
		RoundNumber: e.State.RoundNumber,
		RoundTimer: new(uint64),
	}

	*state.Field = *e.State.Field
	if e.State.RoundTimer != nil {
		atomic.StoreUint64(state.RoundTimer, atomic.LoadUint64(e.State.RoundTimer))
	}
	return state
}

func (e *Engine) updateFieldRound() {
	swampNumber := SwampStartNumber
	for i := 0; i < FieldWidth*FieldHeight; i++ {
		cell := &e.State.Field.Cells[i]
		switch cell.Type {
		case Swamp:
			{
				cell.Type = Water
			}
		case Sand:
			{
				if swampNumber != 0 && rand.Intn(2) == 1 {
					cell.Type = Swamp
					swampNumber--
				}
			}
		}
	}

	for idStr, player := range e.State.Players {
		if player.LoseRound == nil && e.isPlayerDead(idStr) {
			player.LoseRound = &e.State.RoundNumber
			e.State.Players[idStr] = player
			e.GameOver.Store(true)
		}
	}
}

func (e *Engine) updateState(actions *[]*game.Action) {
	e.UpdateM.Lock()
	defer e.UpdateM.Unlock()
	e.StateDiff = initStateDiff()

	for _, action := range *actions {
		switch action.Type {
		case game.InitPlayers:
			{
				e.setGameStart()
				e.initPlayers(action)
				go e.Transport.SendOut(&game.Action{
					Type:    game.SetState,
					Payload: *e.copyState(),
				})
				return
			}
		case game.InitPlayerReady:
			{
				e.initPlayerReady(action)
			}
		case game.InitPlayerMove:
			{
				e.movePlayer(action)

				if e.GameOver.Load().(bool) {
					e.Transport.SendOut(&game.Action{
						Type:    game.SetStateDiff,
						Payload: *e.StateDiff,
					})

					e.Transport.SendOut(&game.Action{
						Type: game.SetGameOver,
					})
					return
				}
			}
		case game.SetRoundStart:
			{
				e.setRoundStart()
				e.RoundRunning.Store(true)
			}
		case game.SetRoundTime:
			{
				e.setRoundTime()
			}
		case game.SetRoundStop:
			{
				e.RoundRunning.Store(false)
				e.setRoundStop()
				go func(state game.State) {
					e.Transport.SendOut(&game.Action{
						Type:    game.SetStateDiff,
						Payload: state,
					})
					e.Transport.SendOut(&game.Action{Type: game.SetRoundStop})
				}(*e.StateDiff)

				e.ReceivedActions <- &game.Action{
					Type: game.SetFieldRound,
				}
				return
			}
		case game.SetFieldRound:
			{
				e.updateFieldRound()
				e.Transport.SendOut(&game.Action{
					Type:    game.SetState,
					Payload: *e.copyState(),
				})

				if e.GameOver.Load().(bool) {
					e.Transport.SendOut(&game.Action{
						Type: game.SetGameOver,
					})
				}
				return
			}
		case game.InitEngineStop:
			{
				e.GameOver.Store(true)
				go e.Transport.SendOut(&game.Action{
					Type: game.SetEngineStop,
				})
				return
			}
		}
	}

	if !e.StateDiff.Empty() {
		go e.Transport.SendOut(&game.Action{
			Type:    game.SetStateDiff,
			Payload: *e.StateDiff,
		})
	}
}

func (e *Engine) run() {
	go e.collectActions()
	e.Ticker = time.NewTicker(TickerDuration * time.Millisecond)

	for {
		select {
		case <-e.Ticker.C:
			{
				if e.GameOver.Load().(bool) {
					close(e.ReceivedActions)
					fmt.Println("Engine stopped")
					return
				}

				e.ProcessM.Lock()
				if len(e.ProcessActions) == 0 {
					e.ProcessM.Unlock()
					continue
				}
				actions := make([]*game.Action, len(e.ProcessActions))
				copy(actions, e.ProcessActions)
				go e.updateState(&actions)
				e.ProcessActions = make([]*game.Action, 0, 100)
				e.ProcessM.Unlock()
			}
		}
	}
}

func (e *Engine) collectActions() {
	for {
		select {
		case action, hasMore := <-e.ReceivedActions:
			{
				if !hasMore {
					fmt.Println("Stop collect actions")
					e.Transport.SendOut(&game.Action{
						Type: game.SetEngineStop,
					})
					return
				}
				e.ProcessM.Lock()
				e.ProcessActions = append(e.ProcessActions, action)
				e.ProcessM.Unlock()
				if action.Type == game.InitEngineStop {
					close(e.ReceivedActions)
					fmt.Println("Stop collect actions")
					return
				}
			}
		}
	}
}

func InitEngine(callback func(action *game.Action)) func(action interface{}) {
	engine := &Engine{
		Transport: &Transport{
			OuterReceiver: callback,
		},
		UpdateM: &sync.Mutex{},
		ProcessM:        &sync.Mutex{},
		ReceivedActions: make(chan *game.Action, 100),
		ProcessActions:  make([]*game.Action, 0, 10),
		RoundRunning: &atomic.Value{},
		GameOver: &atomic.Value{},
	}

	engine.RoundRunning.Store(false)
	engine.GameOver.Store(false)

	engine.Transport.InnerReceiver = func(action interface{}) {
		if isRoundRunning := engine.RoundRunning.Load().(bool); isRoundRunning {
			switch action.(*game.Action).Type {
			case game.InitPlayers, game.InitPlayerReady:
				{
					return
				}
			}
		} else {
			switch action.(*game.Action).Type {
			case game.InitPlayerMove:
				{
					return
				}
			}
		}
		engine.ReceivedActions <- action.(*game.Action)
	}

	go engine.run()

	return engine.Transport.InnerReceiver
}

func InitEngineJS(callback func(actionType, payload string)) func(action interface{}) {
	engine := &Engine{
		Transport: &Transport{
			OuterReceiver: func(action *game.Action) {
				if action.Payload == nil {
					callback(action.Type, "")
					return
				}

				payload, _ := json.Marshal(action.Payload)
				callback(action.Type, string(payload))
			},
		},
		UpdateM: &sync.Mutex{},
		ProcessM:        &sync.Mutex{},
		ReceivedActions: make(chan *game.Action, 100),
		ProcessActions:  make([]*game.Action, 0, 10),
		RoundRunning: &atomic.Value{},
		GameOver: &atomic.Value{},
	}

	engine.RoundRunning.Store(false)
	engine.GameOver.Store(false)

	engine.Transport.InnerReceiver = func(action interface{}) {
		isRoundRunning := engine.RoundRunning.Load().(bool)

		raw := game.ActionRaw{}
		json.Unmarshal([]byte(action.(string)), &raw)

		act := &game.Action{
			Type: raw.Type,
		}

		var payload interface{}
		switch act.Type {
		case game.InitPlayers:
			{
				if isRoundRunning {
					return
				}
				payload = &game.InitPlayersPayload{}
			}
		case game.InitPlayerMove:
			{
				if !isRoundRunning {
					return
				}
				payload = &game.InitPlayerMovePayload{}
			}
		case game.InitPlayerReady:
			{
				if isRoundRunning {
					return
				}
				payload = &game.InitPlayerReadyPayload{}
			}
		default:
			return
		}

		json.Unmarshal([]byte(raw.Payload), payload)
		act.Payload = payload

		engine.ReceivedActions <- act
	}

	go engine.run()

	return engine.Transport.InnerReceiver
}
