package types // import "github.com/daghack/battleground/game/logic"

import (
	"encoding/json"
	"fmt"
)

type GameUpdater func(*Game) error

func (g *GameManager) GameUpdates(gameId string, updates ...GameUpdater) error {
	gamestate, err := g.FetchGame(gameId)
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

func (g *GameManager) CreatePlayer(nickname, passkey string) (string, error) {
	player := &Player{
		Username : nickname,
		Passkey : passkey,
	}
	return g.dbh.CreatePlayer(player)
}

func (g *GameManager) PlayerLogin(nickname, passkey string) (string, error) {
	// Obviously, in a real world application, this behavior is blatantly incorrect,
	// as there is nothing that stops a user from looking at gamestate, then copying another
	// player's id. Even so, it is *good enough* as long as nobody ever plays it.
	// Consider this an area that needs to be fixed the second there are more than 2
	// players in the whole world.
	return g.dbh.VerifyPlayer(nickname, passkey)
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
	players, err := g.dbh.ActivePlayersInGame(gameId)
	if err != nil {
		return err
	}
	if len(players) >= 2 {
		return fmt.Errorf("Game already has the max number of players joined.")
	}
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
	if gamestate.UnitMap == nil {
		gamestate.UnitMap = map[Location]Unit{}
	}
	if gamestate.PlayerStatus == nil {
		gamestate.PlayerStatus = map[string]int{}
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
				loc.Y = 7
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
