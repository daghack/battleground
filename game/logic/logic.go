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

func (gs *GameState) TakeTurn(moveFrom, moveTo Location, orientTowards Orientation, attack bool) error {
	if _, ok := gs.BoardState[moveFrom]; !ok {
		return fmt.Errorf("No Unit At That Position")
	} else if _, ok := gs.BoardState[moveTo]; ok {
		return fmt.Errorf("Another Unit At That Position")
	}
	unit := gs.BoardState[moveFrom]
	unit.Orientation = orientTowards
	delete(gs.BoardState, moveFrom)
	gs.BoardState[moveTo] = unit
	return nil
}
