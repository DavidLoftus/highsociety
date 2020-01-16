package http

import (
	"fmt"
	"github.com/teris-io/shortid"
	"sync"
	"time"
)

type gameLobby struct {
	id      string
	players []*Player
	mut     sync.RWMutex
}

func (lobby *gameLobby) AddPlayer(player *Player) error {
	// TODO: careful we aren't locking when game has started
	lobby.mut.Lock()
	defer lobby.mut.Unlock()

	if len(lobby.players) >= 6 {
		return fmt.Errorf("cannot join full lobby")
	}
	lobby.players = append(lobby.players, player)

	return nil
}

func (lobby *gameLobby) StartGame() {
}

type gameLobbyStore struct {
	lobbies map[string]*gameLobby
	mut     sync.RWMutex
	sid     *shortid.Shortid
}

func (store *gameLobbyStore) NewLobby() (*gameLobby, error) {
	id, err := store.sid.Generate()
	if err != nil {
		return nil, err
	}

	lobby := &gameLobby{
		id: id,
	}

	store.mut.Lock()
	defer store.mut.Unlock()

	if store.lobbies == nil {
		store.lobbies = make(map[string]*gameLobby)
	}
	store.lobbies[id] = lobby

	return lobby, nil
}

func (store *gameLobbyStore) GetLobby(id string) *gameLobby {
	store.mut.RLock()
	defer store.mut.RUnlock()
	return store.lobbies[id]
}

func NewGameLobbyStore() *gameLobbyStore {
	sid := shortid.MustNew(1, shortid.DefaultABC, uint64(time.Now().UnixNano()))
	return &gameLobbyStore{
		sid: sid,
	}
}

var globalGameLobbyStore *gameLobbyStore = NewGameLobbyStore()
