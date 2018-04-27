package types // import "github.com/daghack/battleground/game/logic"

import (
	"fmt"
)

func NewGame(playerId Id) *GameState {
	return &GameState{
		Players : map[Id]*PlayerState{
			playerId : &PlayerState{Id : playerId},
		},
		CurrentPlayer : playerId,
	}
}

func (gs *GameState) AddPlayer (playerId Id) error {
	if _, ok := gs.Players[playerId]; ok {
		return fmt.Errorf("Player Already Exists")
	}
	gs.Players[playerId] = &PlayerState{
		Id : playerId,
	}
	return nil
}

func (gs *GameState) ReadyPlayer(playerId Id, units []UnitType) error {
	if _, ok := gs.Players[playerId]; ok {
		gs.Players[playerId].Ready = true
		gs.Players[playerId].Field = units
		return nil
	}
	return fmt.Errorf("Player Has Not Joined Game")
}
