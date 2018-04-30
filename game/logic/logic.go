package types // import "github.com/daghack/battleground/game/logic"

import (
	"encoding/json"
	"fmt"
)

type GameUpdater func(*Game) error

func (g *GameManager) GameUpdates(gameId string, updates ...GameUpdater) error {
	activeGame, err := g.dbh.FetchGame(gameId)
	if err != nil {
		return err
	}
	gamestate := &Game{}
	err = json.Unmarshal(activeGame.GameState, gamestate)
	if err != nil {
		return err
	}
	for _, v := range updates {
		err = v(gamestate)
		if err != nil {
			return err
		}
	}
	return g.dbh.UpdateGame(gameId, gamestate)
}

type GameManager struct {
	dbh *DBHandler
}

func NewGameManager(dbstring, loadfile string) (*GameManager, error) {
	dbh, err := NewDBHandler(dbstring, loadfile)
	if err != nil {
		return nil, err
	}
	return &GameManager {
		dbh : dbh,
	}, nil
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
	err := g.dbh.JoinGame(gameId, playerId)
	if err != nil {
		return err
	}
	return g.GameUpdates(gameId, updatePlayerStatus(playerId, STATUS_JOINED))
}

func (g *GameManager) ReadyPlayer(gameId, playerId string, unitTypes []string) error {
	return g.GameUpdates(gameId, placeUnits(playerId, unitTypes), updatePlayerStatus(playerId, STATUS_READY))
}

func (g *GameManager) MoveUnit(gameId, playerId string, srcLocation, dstLocation Location, orientation int) error {
	return g.GameUpdates(gameId, moveUnit(playerId, srcLocation, dstLocation, orientation))
}

func (g *GameManager) FetchGame(gameId string) (*Game, error) {
	activeGame, err := g.dbh.FetchGame(gameId)
	if err != nil {
		return nil, err
	}
	gamestate := &Game{}
	err = json.Unmarshal(activeGame.GameState, gamestate)
	if err != nil {
		return nil, err
	}
	return gamestate, nil
}

func updatePlayerStatus(playerId string, STATUS int) GameUpdater {
	return func(gamestate *Game) error {
		gamestate.PlayerStatus[playerId] = STATUS
		return nil
	}
}

func placeUnits(playerId string, unitTypes []string) GameUpdater {
	return func(gamestate *Game) error {
		checkloc := Location{X:0, Y:0}
		for i, unitType := range unitTypes {
			loc := Location{X:i}
			if _, ok := gamestate.UnitMap[checkloc]; ok {
				loc.Y = 8
			}
			gamestate.UnitMap[loc] = Unit{
				Orientation : ORIENT_NORTH,
				UnitType : unitType,
				PlayerId : playerId,
			}
		}
		return nil
	}
}

func moveUnit(playerId string, srcLocation, dstLocation Location, orientation int) GameUpdater {
	return func(gamestate *Game) error {
		if _, ok := gamestate.UnitMap[srcLocation]; !ok {
			return fmt.Errorf("No unit exists at src location.")
		} else if _, ok := gamestate.UnitMap[dstLocation]; ok {
			return fmt.Errorf("Unit already exists at dst location.")
		}
		unit := gamestate.UnitMap[srcLocation]
		if unit.PlayerId != playerId {
			return fmt.Errorf("This unit does not belong to you.")
		}
		unit.Orientation = orientation
		delete(gamestate.UnitMap, srcLocation)
		gamestate.UnitMap[dstLocation] = unit
		return nil
	}
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
