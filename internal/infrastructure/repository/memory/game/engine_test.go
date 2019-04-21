package game

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/domain/game"
	"github.com/stretchr/testify/assert"
)

var testCaseControlRemainingRoundTime = []struct {
	name     string
	gameOver bool
}{
	{
		name:     "Test with gameOver flag",
		gameOver: true,
	},

	{
		name:     "Test without gameOver flag",
		gameOver: false,
	},
}

var testCaseStopRound = []struct {
	name     string
	gameOver bool
}{
	{
		name:     "Test with gameOver flag",
		gameOver: true,
	},

	{
		name:     "Test without gameOver flag",
		gameOver: false,
	},
}

var testCaseInitPlayers = []struct {
	name              string
	playersInitialize bool
}{
	{
		name:              "Test with already initialized players",
		playersInitialize: true,
	},

	{
		name:              "Test with not initialized players",
		playersInitialize: false,
	},
}

var testCaseMovePlayer = []struct {
	name      string
	playerId  uint64
	direction string
	death     bool
}{
	{
		name:      "Test move right",
		playerId:  1,
		direction: Right,
		death:     true,
	},

	{
		name:      "Test move down",
		playerId:  1,
		direction: Down,
		death:     true,
	},
	{
		name:      "Test move left",
		playerId:  1,
		direction: Left,
		death:     false,
	},
	{
		name:      "Test move up",
		playerId:  1,
		direction: Up,
		death:     false,
	},
}

func TestInitState(t *testing.T) {
	expected := &game.State{
		Players:     make(map[string]game.Player),
		ActiveItems: sync.Map{},
		RoundNumber: 0,
	}

	state := initState()
	assert.NotEmpty(t, state.Field,
		"Returns state with empty field")

	expected.Field = state.Field
	assert.Equal(t, expected, state,
		"Creates incorrect state")
}

func TestInitField(t *testing.T) {
	field := initField()
	assert.NotEmpty(t, field, "Return empty result")
}

func TestSetGameState(t *testing.T) {
	engine := &Engine{}

	expected := &game.State{
		Players:     make(map[string]game.Player),
		ActiveItems: sync.Map{},
		RoundNumber: 0,
	}

	engine.setGameState()
	assert.NotEmpty(t, engine.State.Field,
		"Set empty state field")

	expected.Field = engine.State.Field
	assert.Equal(t, expected, engine.State,
		"Sets incorrect state")
}

func TestSetRoundTime(t *testing.T) {
	roundTimer := uint64(5)

	engine := &Engine{
		State: &game.State{
			RoundTimer: new(uint64),
		},
		StateDiff: &game.State{
			RoundTimer: new(uint64),
		},
	}

	*engine.State.RoundTimer = roundTimer
	*engine.StateDiff.RoundTimer = roundTimer

	engine.setRoundTime()

	assert.Equal(t, roundTimer-1, *engine.State.RoundTimer,
		"Doesn't decrease state timer")

	assert.Equal(t, roundTimer-1, *engine.StateDiff.RoundTimer,
		"Doesn't save new roundTimer value to stateDiff")
}

func TestControlRemainingRoundTime(t *testing.T) {
	for _, testCase := range testCaseControlRemainingRoundTime {
		t.Run(testCase.name, func(t *testing.T) {
			engine := &Engine{
				ReceivedActions: make(chan *game.Action, 100),
				GameOver:        &atomic.Value{},
				State: &game.State{
					RoundTimer: new(uint64),
				},
			}

			*engine.State.RoundTimer = uint64(3)

			engine.GameOver.Store(false)
			if testCase.gameOver {
				engine.GameOver.Store(true)
			}

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				for range engine.ReceivedActions {
					atomic.StoreUint64(engine.State.RoundTimer, atomic.LoadUint64(engine.State.RoundTimer)-1)
				}
				wg.Done()
			}()

			engine.controlRemainingRoundTime()
			close(engine.ReceivedActions)
			wg.Wait()

			if testCase.gameOver {
				assert.Equal(t, uint64(3), *engine.State.RoundTimer,
					"Updates timer with gameOver flag")
				return
			}

			assert.Equal(t, uint64(1), *engine.State.RoundTimer,
				"Updates timer with gameOver flag")
		})
	}
}

func TestStopRound(t *testing.T) {
	for _, testCase := range testCaseStopRound {
		t.Run(testCase.name, func(t *testing.T) {
			engine := &Engine{
				ReceivedActions: make(chan *game.Action, 100),
				GameOver:        &atomic.Value{},
			}

			engine.GameOver.Store(false)
			if testCase.gameOver {
				engine.GameOver.Store(true)
			}

			engine.stopRound()
			close(engine.ReceivedActions)
			_, hasMore := <-engine.ReceivedActions

			assert.Equal(t, testCase.gameOver, !hasMore,
				"Incorrect work with channel")
		})
	}
}

func TestInitPlayers(t *testing.T) {
	for _, testCase := range testCaseInitPlayers {
		t.Run(testCase.name, func(t *testing.T) {
			engine := &Engine{
				ReceivedActions: make(chan *game.Action, 100),
				State: &game.State{
					Players: make(map[string]game.Player),
				},
			}

			if testCase.playersInitialize {
				engine.State.Players["100"] = game.Player{
					Id: 100,
				}
			}

			action := &game.Action{
				Payload: &game.InitPlayersPayload{
					PlayersId: []uint64{1},
				},
			}

			engine.initPlayers(action)

			idStr := strconv.FormatUint(uint64(1), 10)
			_, exists := engine.State.Players[idStr]
			assert.Equal(t, testCase.playersInitialize, !exists,
				"Incorrect initialization of players")
		})
	}
}

func TestInitPlayerReady(t *testing.T) {
	engine := &Engine{
		ReceivedActions: make(chan *game.Action, 1),
		State: &game.State{
			Players: make(map[string]game.Player),
		},
	}

	engine.State.Players["1"] = game.Player{
		Id: 1,
	}

	engine.initPlayerReady(&game.Action{
		Payload: &game.InitPlayerReadyPayload{
			PlayerId: 1,
		},
	})
	close(engine.ReceivedActions)

	action := <-engine.ReceivedActions
	assert.Equal(t, true, engine.State.Players["1"].Ready,
		"Doesn't correctly set ready status")
	assert.Equal(t, game.SetRoundStart, action.Type,
		"Doesn't set up SetRoundStart action")
}

func TestMovePlayer(t *testing.T) {
	field := &game.Field{
		Width:  2,
		Height: 2,
		Cells:  make([]game.Cell, 0, 4),
	}

	for i := 0; i < field.Height; i++ {
		for j := 0; j < field.Width; j++ {
			cell := game.Cell{
				Row:  i,
				Col:  j,
				Type: Water,
			}
			if i == 0 && j == 0 {
				cell.Type = Sand
			}
			field.Cells = append(field.Cells, cell)
		}
	}

	engine := &Engine{
		State: &game.State{
			Field:       field,
			RoundNumber: 1,
		},
		StateDiff: initStateDiff(),
	}

	for _, testCase := range testCaseMovePlayer {
		t.Run(testCase.name, func(t *testing.T) {
			engine.State.Players = make(map[string]game.Player)
			playerIdStr := strconv.FormatUint(testCase.playerId, 10)
			engine.State.Players[playerIdStr] = game.Player{
				Id: testCase.playerId,
				X:  CellSize - 1,
				Y:  CellSize - 1,
			}
			engine.GameOver = &atomic.Value{}
			engine.GameOver.Store(false)

			engine.movePlayer(&game.Action{
				Payload: &game.InitPlayerMovePayload{
					PlayerId: testCase.playerId,
					Move:     testCase.direction,
				},
			})

			assert.Equal(t, testCase.death, engine.GameOver.Load(),
				"Incorrect death process")
		})
	}
}

func TestCopyState(t *testing.T) {
	engine := &Engine{
		State: &game.State{
			Field:      initField(),
			RoundTimer: new(uint64),
		},
	}

	*engine.State.RoundTimer = 5

	assert.Equal(t, engine.State, engine.copyState(),
		"Doesn't return copy of state")
}

func TestUpdateFieldRound(t *testing.T) {
	engine := &Engine{
		State: initState(),
	}

	cells := make([]game.Cell, 0, len(engine.State.Field.Cells))
	copy(cells, engine.State.Field.Cells)

	engine.updateFieldRound()
	assert.NotEqual(t, cells, engine.State.Field.Cells,
		"Doesn't update state field")
}

func TestEngineInitPlayersAction(t *testing.T) {
	actions := make(chan *game.Action, 10)
	send := InitEngine(func(action *game.Action) {
		actions <- action
	})

	send(&game.Action{
		Type: game.InitPlayers,
		Payload: &game.InitPlayersPayload{
			PlayersId: []uint64{1},
		},
	})

	time.Sleep(time.Second)

	send(&game.Action{
		Type: game.InitEngineStop,
	})

	action := <-actions
	assert.Equal(t, game.SetState, action.Type,
		"Doesn't set state on correct players init")

	action = <-actions
	assert.Equal(t, game.SetGameOver, action.Type,
		"Doesn't correctly stop engine")

	action = <-actions
	assert.Equal(t, game.SetEngineStop, action.Type,
		"Doesn't correctly stop engine")
	close(actions)
}

func TestInitEngineJS(t *testing.T) {
	InitEngineJS(func(actionType, payload string) {})
}

func waitForActionWithTimeout(t *testing.T, actionType string, duration time.Duration, channel chan *game.Action) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	select {
	case action := <-channel:
		{
			assert.Equal(t, actionType, action.Type,
				"Doesn't return expected action")
		}
	case <-ticker.C:
		{
			assert.Fail(t, "Returning action timeout")
		}
	}
}

func TestUpdateStateMovePlayer(t *testing.T) {
	received := make(chan *game.Action, 10)
	process := make(chan *game.Action, 10)

	engine := &Engine{
		Transport: &Transport{
			OuterReceiver: func(action *game.Action) {
				received <- action
			},
		},
		State: &game.State{
			Field:       initField(),
			Players:     make(map[string]game.Player),
			RoundTimer:  new(uint64),
			RoundNumber: 1,
		},
		ReceivedActions: process,
		GameOver:        &atomic.Value{},
		UpdateM:         &sync.Mutex{},
	}
	engine.State.Players["1"] = game.Player{
		Id: 1,
	}
	engine.State.Field.Cells[0].Type = Sand
	engine.GameOver.Store(false)

	move := &game.Action{
		Type: game.InitPlayerMove,
		Payload: &game.InitPlayerMovePayload{
			PlayerId: 1,
			Move:     Right,
		},
	}
	updateField := &game.Action{
		Type: game.SetFieldRound,
	}

	initEngineStop := &game.Action{
		Type: game.InitEngineStop,
	}

	engine.updateState(&[]*game.Action{move})
	waitForActionWithTimeout(t, game.SetStateDiff, time.Second, received)

	engine.updateState(&[]*game.Action{updateField})
	waitForActionWithTimeout(t, game.SetState, time.Second, received)

	engine.GameOver.Store(true)
	engine.updateState(&[]*game.Action{move})
	waitForActionWithTimeout(t, game.SetStateDiff, time.Second, received)
	waitForActionWithTimeout(t, game.SetGameOver, time.Second, received)

	engine.updateState(&[]*game.Action{initEngineStop})
	waitForActionWithTimeout(t, game.SetEngineStop, time.Second, received)

	close(received)
}
