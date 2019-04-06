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
		Players: make(chan *game.User, 100),
		Closed:  make(chan *game.Room, 100),
		Log:     log,
	}
}

type Repo struct {
	Rooms  *sync.Map
	Total  int
	TotalM *sync.Mutex

	Playing *sync.Map
	Players chan *game.User
	Closed  chan *game.Room

	Log *zap.Logger
}

func (r *Repo) Run() {
	for {
		select {
		case player := <-r.Players:
			{
				r.Log.Info("User start searching game",
					zap.Uint64("user_id", player.Info.ID),
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

func (r *Repo) addUser(player *game.User) {
	if _, playing := r.Playing.Load(player.Info.ID); playing {
		message := &game.Action{
			Type:    "SET_ALREADY_PLAYING",
			Payload: "Already playing",
		}
		player.Messages <- message

		r.Log.Info("Already playing user trying to start new game",
			zap.Uint64("user_id", player.Info.ID),
			zap.String("session_id", player.SessionID))

		return
	}

	r.Playing.Store(player.Info.ID, nil)

	message := &game.Action{
		Type:    "SET_SEARCHING_OPPONENT",
		Payload: "Searching opponent",
	}
	player.Messages <- message

	if err := r.findRoom(player); err != nil {
		message := &game.Action{
			Type:    "SET_NOT_FIND_OPPONENT",
			Payload: "Can't find opponent",
		}
		player.Messages <- message

		r.Log.Warn("User couldn't find opponent",
			zap.Uint64("user_id", player.Info.ID),
			zap.String("session_id", player.SessionID))
	}

	if player.Room.Total == 2 {
		r.Log.Info("Start game at room",
			zap.String("room_id", player.Room.ID),
		)

		sendInto := InitEngine(player.Room.ActionCallback)

		go player.Room.Run(sendInto)
	}
}

func (r *Repo) findRoom(player *game.User) error {
	var result *game.Room

	r.Rooms.Range(func(key, value interface{}) bool {
		received := value.(*game.Room)
		if received.Total < 2 {
			result = received
		}
		return true
	})

	if result != nil {
		result.Users.Store(player.Info.ID, player)
		result.Total++
		player.Room = result

		r.Log.Info("User connect to room",
			zap.Uint64("user_id", player.Info.ID),
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

	result.Users.Store(player.Info.ID, player)
	result.Total++
	r.Rooms.Store(result.ID, result)
	player.Room = result

	r.Log.Info("User create room",
		zap.Uint64("user_id", player.Info.ID),
		zap.String("session_id", player.SessionID),
		zap.String("room_id", player.Room.ID))

	return nil
}

func (r *Repo) PlayersChan() chan *game.User {
	return r.Players
}
