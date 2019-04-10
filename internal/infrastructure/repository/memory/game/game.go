package game

import (
	"go.uber.org/zap"
	"sadislands/internal/domain/game"
	"sync"
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

				go player.ListenAndSend(r.Log)
				go r.addUser(player)
			}
		case room := <-r.Closed:
			{
				room.Users.Range(func(key, value interface{}) bool {
					r.Playing.Delete(key)
					return true
				})
				r.Rooms.Delete(room.ID)
				r.TotalM.Lock()
				r.Total--
				r.TotalM.Unlock()
				r.Log.Info("Room closed",
					zap.String("room_id", room.ID),
				)
			}
		}
	}
}

func (r *Repo) addUser(player *game.User) {
	if _, playing := r.Playing.Load(player.Info.ID); playing {
		message := &game.Action{
			Type:    game.SetAlreadyPlaying,
		}
		player.Messages <- message

		r.Log.Info("Already playing user trying to start new game",
			zap.Uint64("user_id", player.Info.ID),
			zap.String("session_id", player.SessionID))

		return
	}

	r.Playing.Store(player.Info.ID, nil)

	message := &game.Action{
		Type:    game.SetOpponentSearch,
	}
	player.Messages <- message

	if err := r.findRoom(player); err != nil {
		message := &game.Action{
			Type:    game.SetOpponentNotFound,
		}
		player.Messages <- message

		r.Log.Warn("User couldn't find opponent",
			zap.Uint64("user_id", player.Info.ID),
			zap.String("session_id", player.SessionID))
	}

	if player.Room.Total.Load() == 2 {
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
		if received.Total.Load() < 2 {
			result = received
		}
		return true
	})

	if result != nil {
		result.Users.Store(player.Info.ID, player)
		result.Total.Add(1)
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
	result.Total.Add(1)
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
