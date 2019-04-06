package game

import (
	"fmt"
	"sadislands/internal/domain/game"
	"strconv"
	"sync"
	"time"

	"github.com/cathalgarvey/fmtless/encoding/json"
)

const (
	Sand        = "SAND"
	Water       = "WATER"
	FieldWidth  = 10
	FieldHeight = 10
	Up          = "UP"
	Down        = "DOWN"
	Left        = "LEFT"
	Right       = "RIGHT"
)

type Engine struct {
	Transport       *Transport
	State           *game.State
	StateDiff       *game.State
	ProcessActions  []game.Action
	ProcessM        *sync.Mutex
	ReceivedActions chan game.Action
	Ticker          *time.Ticker
	Timer           *time.Timer
}

func initState() *game.State {
	return &game.State{
		Field:       initField(),
		Players:     make(map[string]game.Player),
		ActiveItems: &sync.Map{},
		RoundNumber: 0,
	}
}

func initField() *game.Field {
	field := &game.Field{
		Cells:  make([]game.Cell, 0, FieldWidth*FieldHeight),
		Width:  FieldWidth,
		Height: FieldHeight,
	}

	for i := 0; i < FieldHeight; i++ {
		for j := 0; j < FieldWidth; j++ {
			cell := game.Cell{
				Row:    i,
				Col:    j,
				Type:   Sand,
				HasBox: false,
			}

			field.Cells = append(field.Cells, cell)
		}
	}

	return field
}

func (e *Engine) setGameStart() {
	e.State = initState()
	e.StateDiff = e.State
}

func (e *Engine) setStartRound() {
	e.State.RoundNumber += 1
	time.AfterFunc(5*time.Second, e.stopRound)
	e.State.RoundTimer = 5
	go e.controlRemainingRoundTime()

	e.StateDiff.RoundNumber = e.State.RoundNumber
	e.StateDiff.RoundTimer = e.State.RoundTimer
}

func (e *Engine) setRoundTime() {
	e.State.RoundTimer -= 1
	e.StateDiff.RoundTimer = e.State.RoundTimer
}

func (e *Engine) controlRemainingRoundTime() {
	action := game.Action{
		Type: game.SetRoundTime,
	}

	t := time.NewTicker(1 * time.Second)
	for {
		<-t.C

		if e.State.RoundTimer == 2 {
			e.ReceivedActions <- action
			t.Stop()
			return
		}

		e.ReceivedActions <- action
	}
}

func (e *Engine) stopRound() {
	e.ReceivedActions <- game.Action{
		Type: game.SetRoundStop,
	}
}

func (e *Engine) initPlayers(action *game.Action) {
	payload := &game.InitPlayersPayload{}
	json.Unmarshal([]byte(action.Payload), payload)
	for _, id := range payload.PlayersId {
		idStr := strconv.FormatUint(id, 10)

		e.State.Players[idStr] = game.Player{
			Id:    id,
			Items: make(map[string]uint64),
		}
	}
	e.StateDiff.Players = e.State.Players
}

func (e *Engine) movePlayer(action *game.Action) {
	payload := &game.InitPlayerMovePayload{}
	json.Unmarshal([]byte(action.Payload), payload)

	idStr := strconv.FormatUint(payload.PlayerId, 10)
	player, exists := e.State.Players[idStr]

	if !exists {
		return
	}

	switch payload.Move {
	case Right:
		if player.X == FieldWidth {
			return
		}

		player.X += 1
	case Up:
		if player.Y == 0 {
			return
		}

		player.Y -= 1
	case Left:
		if player.X == 0 {
			return
		}

		player.X -= 1
	case Down:
		if player.Y == FieldHeight {
			return
		}

		player.Y += 1
	}

	e.State.Players[idStr] = player
	e.StateDiff.Players[idStr] = player
}

func (e *Engine) sendStateDiff() {
	payload, _ := json.Marshal(e.StateDiff)
	action, _ := json.Marshal(&game.Action{
		Type:    game.SetState,
		Payload: string(payload),
	})
	go e.Transport.SendOut(game.SetState, string(action))
}

func (e *Engine) setRoundStop() {
	e.State.RoundTimer = 0
	e.StateDiff.RoundTimer = e.State.RoundTimer
}

func actionJSON(actionType string, payload interface{}) string {
	payloadRaw, _ := json.Marshal(payload)
	actionRaw, _ := json.Marshal(&game.Action{
		Type:    actionType,
		Payload: string(payloadRaw),
	})
	return string(actionRaw)
}

func initStateDiff() *game.State {
	return &game.State{
		Players: make(map[string]game.Player),
	}
}

func (e *Engine) updateState(actions *[]game.Action) {
	e.StateDiff = initStateDiff()

	for _, action := range *actions {
		switch action.Type {
		case game.InitPlayers:
			{
				e.setGameStart()
				e.initPlayers(&action)
				e.setStartRound()
			}
		case game.InitPlayerMove:
			{
				e.movePlayer(&action)
			}
		case game.SetRoundTime:
			{
				e.setRoundTime()
			}
		case game.SetRoundStop:
			{
				e.setRoundStop()
				e.sendStateDiff()
				fmt.Println(e.StateDiff)
				go e.Transport.SendOut(game.SetRoundStop, actionJSON(game.SetRoundStop, nil))
				return
			}
		}
	}

	e.sendStateDiff()
}

func (e *Engine) run() {
	go e.collectActions()
	e.Ticker = time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-e.Ticker.C:
			{
				if len(e.ProcessActions) == 0 {
					fmt.Println("Waiting for Actions...")
					continue
				}

				fmt.Println("Process", e.ProcessActions)
				e.ProcessM.Lock()
				actions := make([]game.Action, len(e.ProcessActions))
				copy(actions, e.ProcessActions)
				go e.updateState(&actions)
				e.ProcessActions = make([]game.Action, 0, 100)
				e.ProcessM.Unlock()
			}
		}
	}
}

func (e *Engine) collectActions() {
	for {
		select {
		case action := <-e.ReceivedActions:
			{
				e.ProcessM.Lock()
				e.ProcessActions = append(e.ProcessActions, action)
				e.ProcessM.Unlock()
			}
		}
	}
}

func (e *Engine) InitEngine(callback func(actionType, payload string)) func(action string) {
	e.Transport = &Transport{}
	e.ProcessM = &sync.Mutex{}
	e.ReceivedActions = make(chan game.Action, 100)
	e.ProcessActions = make([]game.Action, 0, 10)

	e.Transport.OuterReceiver = callback
	e.Transport.InnerReceiver = func(action string) {
		act := game.Action{}
		json.Unmarshal([]byte(action), &act)

		e.ReceivedActions <- act
	}

	go e.run()

	return e.Transport.InnerReceiver
}
