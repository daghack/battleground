package types // import "github.com/daghack/battleground/game/logic"

import (
	"encoding/json"
	"fmt"
)

type GameManager struct {
	dbh *DBHandler
}

func NewGameManager() (*GameManager, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (g *GameManager) CreateGame(playerId string) (string, error) {
	gameId, err := g.dbh.CreateGame(8, 8)
	if err != nil {
		return "", err
	}
	err = g.JoinGame(gameId, playerId)
	return gameId, err
}

func (g *GameManager) JoinGame(gameId, playerId string) error {
	activeGame, err := g.dbh.FetchGame(gameId)
	if err != nil {
		return err
	}
	gamestate := &Game{}
	err = json.Unmarshal(activeGame.GameState, gamestate)
	if err != nil {
		return err
	}
	gamestate.PlayerStatus[playerId] = STATUS_JOINED
	return g.dbh.UpdateGame(gameId, gamestate)
}

/*
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
*/
