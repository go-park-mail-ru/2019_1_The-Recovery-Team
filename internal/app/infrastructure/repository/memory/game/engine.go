package game

import (
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/game"

	"github.com/cathalgarvey/fmtless/encoding/json"
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
	Lifebuoy                = "LIFEBUOY"
	Bomb                    = "BOMB"
	ItemDuration     uint64 = 5
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

	GameStart *atomic.Value
	GameOver  *atomic.Value
	Stopped   *atomic.Value

	ReceivedActions chan *game.Action

	Ticker *time.Ticker
	Timer  *time.Timer
}

// initState creates instance of default state
func initState() *game.State {
	return &game.State{
		Field:       initField(),
		Players:     make(map[string]game.Player),
		ActiveItems: make(map[uint64]game.Item),
		RoundNumber: 0,
	}
}

// initField creates instance of field with random cells
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

	// Shuffle game cell types
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(types), func(i, j int) { types[i], types[j] = types[j], types[i] })

	// Player position shouldn't be water
	types[0] = Sand

	// Initialize game field
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

// setGameState sets up game state
func (e *Engine) setGameState() {
	e.State = initState()
}

// Item action handlers

func (e *Engine) initItemUse(action *game.Action) {
	payload := action.Payload.(*game.InitItemUsePayload)

	idStr := strconv.FormatUint(payload.PlayerId, 10)
	player, exists := e.State.Players[idStr]

	// Player exists or already lost
	if !exists || player.LoseRound != nil {
		return
	}

	if n, exists := player.Items[payload.ItemType]; !exists || n == 0 {
		return
	}

	if _, exists := e.State.ActiveItems[payload.PlayerId]; exists {
		return
	}

	player.Items[payload.ItemType]--
	e.setItemStart(payload.PlayerId, payload.ItemType)
}

func (e *Engine) controlRemainingItemTime(playerId uint64, itemType string) {
	remaining := ItemDuration
	action := &game.Action{
		Type: game.SetItemTime,
		Payload: &game.SetItemPayload{
			PlayerId: playerId,
			ItemType: itemType,
		},
	}

	t := time.NewTicker(1 * time.Second)
	for {
		<-t.C
		if e.GameOver.Load().(bool) {
			return
		}

		// Remaining time will become 1 second
		if remaining == 2 {
			e.ReceivedActions <- action
			t.Stop()
			return
		}

		remaining--
		e.ReceivedActions <- action
	}
}

func (e *Engine) setItemStart(playerId uint64, itemType string) {
	e.State.ActiveItems[playerId] = game.Item{
		Type:     itemType,
		PlayerId: playerId,
		Duration: ItemDuration,
	}
	for key, value := range e.State.ActiveItems {
		e.StateDiff.ActiveItems[key] = value
	}

	time.AfterFunc(time.Duration(ItemDuration)*time.Second, e.stopItem(playerId, itemType))
	go e.controlRemainingItemTime(playerId, itemType)
}

func (e *Engine) setItemTime(action *game.Action) {
	payload := action.Payload.(*game.SetItemPayload)
	item := e.State.ActiveItems[payload.PlayerId]
	item.Duration--
	e.State.ActiveItems[payload.PlayerId] = item

	for key, value := range e.State.ActiveItems {
		e.StateDiff.ActiveItems[key] = value
	}
}

func (e *Engine) setItemStop(action *game.Action) {
	payload := action.Payload.(*game.SetItemPayload)
	delete(e.State.ActiveItems, payload.PlayerId)
	delete(e.StateDiff.ActiveItems, payload.PlayerId)

	// Check player death
	idStr := strconv.Itoa(int(payload.PlayerId))
	if payload.ItemType == Lifebuoy && e.isPlayerDead(idStr) {
		player := e.State.Players[idStr]
		player.LoseRound = &e.State.RoundNumber
		e.GameOver.Store(true)
		e.StateDiff.Players[idStr] = player
	}
}

func (e *Engine) stopItem(playerId uint64, itemType string) func() {
	return func() {
		if e.GameOver.Load().(bool) {
			return
		}

		e.ReceivedActions <- &game.Action{
			Type: game.SetItemStop,
			Payload: &game.SetItemPayload{
				PlayerId: playerId,
				ItemType: itemType,
			},
		}
	}
}

// Round action handlers

// setRoundStart updates round number and starts round timer
func (e *Engine) setRoundStart() {
	e.State.RoundNumber += 1

	// Create setStopRound action on time expire
	time.AfterFunc(time.Duration(RoundDuration)*time.Second, e.stopRound)

	e.State.RoundTimer = new(uint64)
	e.StateDiff.RoundTimer = new(uint64)
	atomic.StoreUint64(e.State.RoundTimer, RoundDuration)

	go e.controlRemainingRoundTime()

	e.StateDiff.RoundNumber = e.State.RoundNumber

	atomic.StoreUint64(e.StateDiff.RoundTimer, atomic.LoadUint64(e.State.RoundTimer))
}

// setRoundTime decreases round timer
func (e *Engine) setRoundTime() {
	e.StateDiff.RoundTimer = new(uint64)
	atomic.StoreUint64(e.State.RoundTimer, atomic.LoadUint64(e.State.RoundTimer)-1)
	atomic.StoreUint64(e.StateDiff.RoundTimer, atomic.LoadUint64(e.State.RoundTimer))
}

// controlRemainingRoundTime controls update of timer
func (e *Engine) controlRemainingRoundTime() {
	action := &game.Action{
		Type: game.SetRoundTime,
	}

	t := time.NewTicker(1 * time.Second)
	for {
		<-t.C
		if e.GameOver.Load().(bool) {
			return
		}

		// Remaining time will become 1 second
		if atomic.LoadUint64(e.State.RoundTimer) == 2 {
			e.ReceivedActions <- action
			t.Stop()
			return
		}

		e.ReceivedActions <- action
	}
}

// stopRound creates SetRoundStop action
func (e *Engine) stopRound() {
	if e.GameOver.Load().(bool) {
		return
	}

	e.ReceivedActions <- &game.Action{
		Type: game.SetRoundStop,
	}
}

// Players action handlers

// initPlayers initialize players in current game
func (e *Engine) initPlayers(action *game.Action) {
	if len(e.State.Players) != 0 {
		return
	}

	payload := action.Payload.(*game.InitPlayersPayload)

	for _, id := range payload.PlayersId {
		idStr := strconv.FormatUint(id, 10)

		player := game.Player{
			Id:    id,
			Items: make(map[string]uint64),
		}
		player.Items[Lifebuoy] = 3
		player.Items[Bomb] = 0
		player.Items[Sand] = 3

		e.State.Players[idStr] = player
	}
}

// initPlayerReady updates player ready status
// and, if all players are ready, creates SetRoundStart action
func (e *Engine) initPlayerReady(action *game.Action) {
	payload := action.Payload.(*game.InitPlayerReadyPayload)

	idStr := strconv.FormatUint(payload.PlayerId, 10)
	player, exists := e.State.Players[idStr]

	if !exists || player.Ready {
		return
	}

	player.Ready = true
	e.State.Players[idStr] = player

	// Check for other players ready status
	for _, player := range e.State.Players {
		if !player.Ready {
			return
		}
	}

	e.GameStart.Store(true)

	e.ReceivedActions <- &game.Action{
		Type: game.SetRoundStart,
	}
}

// movePlayer moves player position and checks player for death
func (e *Engine) movePlayer(action *game.Action) {
	payload := action.Payload.(*game.InitPlayerMovePayload)

	idStr := strconv.FormatUint(payload.PlayerId, 10)
	player, exists := e.State.Players[idStr]

	// Player exists or already lost
	if !exists || player.LoseRound != nil {
		return
	}

	// Process move direction
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
	if item, exists := e.State.ActiveItems[player.Id]; exists && item.Type == Sand {
		e.State.Field.Cells[(player.X/CellSize)+(player.Y/CellSize)*e.State.Field.Width].Type = Sand
		if e.StateDiff.Field == nil {
			e.StateDiff.Field = &game.Field{}

			e.StateDiff.Field.Cells = make([]game.Cell, len(e.State.Field.Cells))
			copy(e.StateDiff.Field.Cells, e.State.Field.Cells)
			e.StateDiff.Field.Height = e.State.Field.Height
			e.StateDiff.Field.Width = e.State.Field.Width

		} else {
			e.StateDiff.Field.Cells[(player.X/CellSize)+(player.Y/CellSize)*e.State.Field.Width].Type = Sand
		}
	} else if e.isPlayerDead(idStr) {
		player.LoseRound = &e.State.RoundNumber
		e.GameOver.Store(true)
	}

	e.StateDiff.Players[idStr] = player
}

// isPlayerDead checks player for death
func (e *Engine) isPlayerDead(id string) bool {
	player := e.State.Players[id]
	if item, exists := e.State.ActiveItems[player.Id]; exists && item.Type == Lifebuoy {
		return false
	}

	return e.State.Field.Cells[(player.X/CellSize)+(player.Y/CellSize)*e.State.Field.Width].Type == Water
}

// initStateDiff creates instance of empty state
func (e *Engine) initStateDiff() {
	e.StateDiff = &game.State{
		Players:     make(map[string]game.Player),
		ActiveItems: make(map[uint64]game.Item),
	}

	if e.State != nil {
		for key, value := range e.State.ActiveItems {
			e.StateDiff.ActiveItems[key] = value
		}
	}
}

// copyState creates copy of current game state
func (e *Engine) copyState() *game.State {
	state := &game.State{
		Field:       &game.Field{},
		Players:     make(map[string]game.Player),
		RoundNumber: e.State.RoundNumber,
		RoundTimer:  new(uint64),
		ActiveItems: make(map[uint64]game.Item),
	}

	for key, value := range e.State.Players {
		state.Players[key] = value
	}
	for key, value := range e.State.ActiveItems {
		state.ActiveItems[key] = value
	}

	state.Field.Cells = make([]game.Cell, len(e.State.Field.Cells))
	copy(state.Field.Cells, e.State.Field.Cells)
	state.Field.Height = e.State.Field.Height
	state.Field.Width = e.State.Field.Width
	if e.State.RoundTimer != nil {
		atomic.StoreUint64(state.RoundTimer, atomic.LoadUint64(e.State.RoundTimer))
	}
	return state
}

// updateFieldRound updates round field after end of round
// and checks all players for death
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

// updateState processes actions and updates state
func (e *Engine) updateState(actions *[]*game.Action) {
	e.UpdateM.Lock()
	defer e.UpdateM.Unlock()

	// Initialize state diff
	e.initStateDiff()

	// Flag for sending empty state
	var forceStateSend bool

	for _, action := range *actions {
		switch action.Type {
		case game.InitPlayers:
			{
				e.setGameState()
				e.initPlayers(action)
				go e.Transport.SendOut(&game.Action{
					Type:    game.SetState,
					Payload: e.copyState(),
				})
				continue
			}
		case game.InitPlayerReady:
			{
				e.initPlayerReady(action)
			}
		case game.InitPlayerMove:
			{
				if e.GameOver.Load().(bool) {
					continue
				}

				e.movePlayer(action)

				// Check game end
				if e.GameOver.Load().(bool) {
					e.Transport.SendOut(&game.Action{
						Type:    game.SetStateDiff,
						Payload: e.StateDiff,
					})

					e.ReceivedActions <- &game.Action{
						Type: game.InitEngineStop,
					}
					return
				}
			}
		case game.SetRoundTime:
			{
				e.setRoundTime()
			}
		case game.SetRoundStart:
			{
				e.setRoundStart()
			}
		case game.SetRoundStop:
			{
				e.updateFieldRound()
				isGameOver := e.GameOver.Load().(bool)
				if !isGameOver {
					e.setRoundStart()
				}
				e.Transport.SendOut(&game.Action{
					Type:    game.SetStateDiff,
					Payload: e.copyState(),
				})

				if isGameOver {
					e.ReceivedActions <- &game.Action{
						Type: game.InitEngineStop,
					}
					return
				}
				e.initStateDiff()
			}
		case game.InitEngineStop:
			{
				e.Transport.SendOut(&game.Action{
					Type: game.SetEngineStop,
				})
				e.Stopped.Store(true)
				e.GameOver.Store(true)
				close(e.ReceivedActions)
				return
			}
		case game.InitItemUse:
			{
				e.initItemUse(action)
			}
		case game.SetItemTime:
			{
				e.setItemTime(action)
			}
		case game.SetItemStop:
			{
				e.setItemStop(action)
				forceStateSend = true

				// Check game end
				if e.GameOver.Load().(bool) {
					e.Transport.SendOut(&game.Action{
						Type:    game.SetStateDiff,
						Payload: e.StateDiff,
					})

					e.ReceivedActions <- &game.Action{
						Type: game.InitEngineStop,
					}
					return
				}
			}
		}
	}

	if forceStateSend || !e.StateDiff.Empty() {
		go e.Transport.SendOut(&game.Action{
			Type:    game.SetStateDiff,
			Payload: e.StateDiff,
		})
	}
}

// run starts game engine
func (e *Engine) run() {
	// Start collecting actions
	go e.collectActions()

	e.Ticker = time.NewTicker(TickerDuration * time.Millisecond)
	for range e.Ticker.C {
		if e.Stopped.Load().(bool) {
			return
		}

		// Start processing actions
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

// collectActions receives actions from channel
// and saves to slice
func (e *Engine) collectActions() {
	for action := range e.ReceivedActions {
		if action.Type == game.InitEngineStop {
			e.Transport.SendOut(&game.Action{
				Type: game.SetGameOver,
			})
		}

		e.ProcessM.Lock()
		e.ProcessActions = append(e.ProcessActions, action)
		e.ProcessM.Unlock()

		// Action for stopping engine
		if action.Type == game.InitEngineStop {
			return
		}
	}
}

// InitEngine creates engine instance
// and returns function to pass actions to engine
func InitEngine(callback func(action *game.Action)) func(action interface{}) {
	engine := &Engine{
		Transport: &Transport{
			OuterReceiver: callback,
		},
		UpdateM:         &sync.Mutex{},
		ProcessM:        &sync.Mutex{},
		ReceivedActions: make(chan *game.Action, 100),
		ProcessActions:  make([]*game.Action, 0, 10),
		GameStart:       &atomic.Value{},
		GameOver:        &atomic.Value{},
		Stopped:         &atomic.Value{},
	}

	// Initialize atomic values with bools inside
	engine.GameStart.Store(false)
	engine.GameOver.Store(false)
	engine.Stopped.Store(false)

	// Initialize innerReceiver logic
	engine.Transport.InnerReceiver = func(action interface{}) {
		if isGameStart := engine.GameStart.Load().(bool); isGameStart {
			switch action.(*game.Action).Type {
			case game.InitPlayers, game.InitPlayerReady:
				{
					return
				}
			}
		} else {
			switch action.(*game.Action).Type {
			case game.InitPlayerMove, game.InitItemUse:
				{
					return
				}
			}
		}

		if engine.Stopped.Load().(bool) {
			return
		}

		engine.ReceivedActions <- action.(*game.Action)
	}

	// Starts engine
	go engine.run()

	return engine.Transport.InnerReceiver
}

// InitEngine creates engine instance for JS
// and returns function to pass actions to engine
func InitEngineJS(callback func(actionType, payload string)) func(action interface{}) {
	engine := &Engine{
		Transport: &Transport{
			OuterReceiver: func(action *game.Action) {
				if action.Payload == nil {
					callback(action.Type, "")
					return
				}

				if action.Type == game.SetStateDiff || action.Type == game.SetState {
					payload := action.Payload.(*game.State)
					activeItems := payload.ActiveItems
					items := make(map[string]game.Item)
					for key, value := range activeItems {
						items[strconv.Itoa(int(key))] = value
					}

					action.Payload = struct {
						Field       *game.Field            `json:"field,omitempty"`
						Players     map[string]game.Player `json:"players,omitempty"`
						ActiveItems map[string]game.Item   `json:"activeItems"`
						RoundNumber int                    `json:"roundNumber,omitempty"`
						RoundTimer  *uint64                `json:"roundTimer,omitempty"`
					}{
						Field:       payload.Field,
						Players:     payload.Players,
						ActiveItems: items,
						RoundNumber: payload.RoundNumber,
						RoundTimer:  payload.RoundTimer,
					}
				}

				payload, _ := json.Marshal(action.Payload)
				callback(action.Type, string(payload))
			},
		},
		UpdateM:         &sync.Mutex{},
		ProcessM:        &sync.Mutex{},
		ReceivedActions: make(chan *game.Action, 100),
		ProcessActions:  make([]*game.Action, 0, 10),
		GameStart:       &atomic.Value{},
		GameOver:        &atomic.Value{},
		Stopped:         &atomic.Value{},
	}

	// Initialize atomic values with bools inside
	engine.GameStart.Store(false)
	engine.GameOver.Store(false)
	engine.Stopped.Store(false)

	// Initialize innerReceiver logic
	engine.Transport.InnerReceiver = func(action interface{}) {
		isRoundRunning := engine.GameStart.Load().(bool)

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
		case game.InitItemUse:
			{
				if !isRoundRunning {
					return
				}
				payload = &game.InitItemUsePayload{}
			}
		default:
			return
		}

		json.Unmarshal([]byte(raw.Payload), payload)
		act.Payload = payload

		engine.ReceivedActions <- act
	}

	// Starts engine
	go engine.run()

	return engine.Transport.InnerReceiver
}
