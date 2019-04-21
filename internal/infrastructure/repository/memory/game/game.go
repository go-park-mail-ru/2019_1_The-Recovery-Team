package game

import (
	"sync"

	"go.uber.org/atomic"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/domain/game"
	"go.uber.org/zap"
)

// NewGameRepo creates new instance of game repository
func NewGameRepo(log *zap.Logger) *Repo {
	repo := &Repo{
		Rooms:   &sync.Map{},
		Total:   &atomic.Uint64{},
		Playing: &sync.Map{},
		Players: make(chan *game.User, 100),
		Closed:  make(chan *game.Room, 100),
		Log:     log,
	}
	repo.Total.Store(0)
	return repo
}

type Repo struct {
	Rooms *sync.Map
	Total *atomic.Uint64

	Playing *sync.Map
	Players chan *game.User
	Closed  chan *game.Room

	Log *zap.Logger
}

// Run starts game handler
func (r *Repo) Run() {
	for {
		select {
		case player := <-r.Players:
			{
				r.Log.Info("User start searching game",
					zap.Uint64("user_id", player.Info.ID),
					zap.String("session_id", player.SessionID))

				// Add user to queue for room
				go r.addUser(player)
			}
		case room := <-r.Closed:
			{
				// Remove users from already playing users
				room.Users.Range(func(key, value interface{}) bool {
					r.Playing.Delete(key)
					r.Log.Info("User can already playing",
						zap.Uint64("user_id", key.(uint64)))
					return true
				})

				// Remove room, update statistics
				r.Rooms.Delete(room.ID)
				r.Total.Sub(1)
				r.Log.Info("Room closed",
					zap.String("room_id", room.ID),
				)
			}
		}
	}
}

func (r *Repo) addUser(player *game.User) {
	// Check user that he is already playing
	if _, playing := r.Playing.Load(player.Info.ID); playing {
		message := &game.Action{
			Type: game.SetAlreadyPlaying,
		}

		if err := player.Conn.WriteJSON(message); err != nil {
			r.Log.Info("Player disconnected")
			return
		}

		r.Log.Info("Already playing user trying to start new game",
			zap.Uint64("user_id", player.Info.ID),
			zap.String("session_id", player.SessionID))

		return
	}

	// Add user for already playing user
	r.Playing.Store(player.Info.ID, nil)

	message := &game.Action{
		Type: game.SetOpponentSearch,
	}

	if err := player.Conn.WriteJSON(message); err != nil {
		r.Log.Info("Player disconnected")
		r.Playing.Delete(player.Info.ID)
		return
	}

	// Searching room
	if err := r.findRoom(player); err != nil {
		message := &game.Action{
			Type: game.SetOpponentNotFound,
		}
		if err := player.Conn.WriteJSON(message); err != nil {
			r.Log.Info("Player disconnected")
			r.Playing.Delete(player.Info.ID)
			return
		}

		r.Log.Warn("User couldn't find opponent",
			zap.Uint64("user_id", player.Info.ID),
			zap.String("session_id", player.SessionID))
	}

	// Activate player listening and sending
	go player.ListenAndSend(player.Room.Log)

	// Start game engine, if room is full
	if player.Room.Total.Load() == 2 {
		r.Log.Info("Start game at room",
			zap.String("room_id", player.Room.ID),
		)

		sendInto := InitEngine(player.Room.ActionCallback)

		// Turn room to game mode
		go player.Room.Run(sendInto)
	}
}

func (r *Repo) findRoom(player *game.User) error {
	var result *game.Room

	// Searching empty room
	r.Rooms.Range(func(key, value interface{}) bool {
		received := value.(*game.Room)
		if received.Total.Load() < 2 {
			result = received
		}
		return true
	})

	// If find room
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

	// Create new room
	result = game.NewRoom(
		r.Log,
		r.Closed,
	)

	r.Total.Add(1)

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

// PlayersChan returns players searching queue
func (r *Repo) PlayersChan() chan *game.User {
	return r.Players
}
