package game

import (
	"sadislands/internal/domain/game"
	"sync"

	"go.uber.org/zap"
)

func NewGameRepo(log *zap.Logger) *Repo {
	return &Repo{
		Rooms:   &sync.Map{},
		TotalM:  &sync.Mutex{},
		Playing: &sync.Map{},
		Players: make(chan *game.Player, 100),
		Closed:  make(chan *game.Room, 100),
		Log:     log,
	}
}

type Repo struct {
	Rooms  *sync.Map
	Total  int
	TotalM *sync.Mutex

	Playing *sync.Map
	Players chan *game.Player
	Closed  chan *game.Room

	Log *zap.Logger
}

func (r *Repo) Run() {
	for {
		select {
		case player := <-r.Players:
			{
				r.Log.Info("User start searching game",
					zap.Uint64("user_id", player.ID),
					zap.String("session_id", player.SessionID))

				go r.addUser(player)
			}
		case room := <-r.Closed:
			{
				r.Rooms.Delete(room.ID)
				r.TotalM.Lock()
				r.Total--
				r.TotalM.Unlock()

			}
		}
	}
}

func (r *Repo) addUser(player *game.Player) {
	if _, playing := r.Playing.Load(player.ID); playing {
		message := &game.Message{
			Status:  400,
			Payload: "Already playing",
		}
		player.Messages <- message

		r.Log.Info("Already playing user trying to start new game",
			zap.Uint64("user_id", player.ID),
			zap.String("session_id", player.SessionID))

		return
	}

	r.Playing.Store(player.ID, nil)

	message := &game.Message{
		Status:  200,
		Payload: "Searching opponent",
	}
	player.Messages <- message

	if err := r.findRoom(player); err != nil {
		message := &game.Message{
			Status:  404,
			Payload: "Can't find opponent",
		}
		player.Messages <- message

		r.Log.Error("User couldn't find opponent",
			zap.Uint64("user_id", player.ID),
			zap.String("session_id", player.SessionID))
	}

	if player.Room.Total == 2 {
		r.Log.Error("Start game at room",
			zap.String("room_id", player.Room.ID),
		)

		go player.Room.Run()
	}
}

func (r *Repo) findRoom(player *game.Player) error {
	var result *game.Room

	r.Rooms.Range(func(key, value interface{}) bool {
		received := value.(*game.Room)
		if received.Total < 2 {
			result = received
		}
		return true
	})

	if result != nil {
		result.Players.Store(player.ID, player)
		result.Total++
		player.Room = result

		r.Log.Error("User connect to room",
			zap.Uint64("user_id", player.ID),
			zap.String("session_id", player.SessionID),
			zap.String("room_id", player.Room.ID))

		return nil
	}

	result = game.NewRoom(
		r.Log,
		r.Closed,
	)

	r.TotalM.Lock()
	r.Total++
	r.TotalM.Unlock()

	result.Players.Store(player.ID, player)
	result.Total++
	r.Rooms.Store(result.ID, result)
	player.Room = result

	r.Log.Error("User create room",
		zap.Uint64("user_id", player.ID),
		zap.String("session_id", player.SessionID),
		zap.String("room_id", player.Room.ID))

	return nil
}

func (r *Repo) PlayersChan() chan *game.Player {
	return r.Players
}
