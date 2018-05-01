package types // import "github.com/daghack/battlegrounds/game/logic"

import (
	"fmt"
	"encoding/json"
)

const (
	ORIENT_NORTH int = iota
	ORIENT_SOUTH
	ORIENT_EAST
	ORIENT_WEST
)

const (
	STATUS_JOINED int = iota
	STATUS_READY  int = iota
)

type Location struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Unit struct {
	Orientation int    `json:"orientation"`
	UnitType    string `json:"unitType"`
	PlayerId    string `json:"playerId"`
}

type UnitMap map[Location]Unit

type Player struct {
	Id       string `db:"id"`
	Username string `db:"username"`
	Passkey  string `db:"passkey"`
}

type ActiveGame struct {
	Id         string `db:"id"`
	BoardSize  int    `db:"board_size"`
	PieceCount int    `db:"piece_count"`
	GameState  []byte `db:"game_state"`
}

type Game struct {
	UnitMap      UnitMap        `json:"unit_map"`
	PlayerStatus map[string]int `json:"player_status"`
	Turn         int            `json:"turn"`
}

type jsonGame struct {
	UnitMap      map[string]Unit `json:"unit_map"`
	PlayerStatus map[string]int `json:"player_status"`
	Turn         int            `json:"turn"`
}

func (game *Game) MarshalJSON() ([]byte, error) {
	fmt.Println("Using the Marshal Function")
	jgame := &jsonGame{
		UnitMap : map[string]Unit{},
		PlayerStatus : game.PlayerStatus,
		Turn : game.Turn,
	}
	for loc, unit := range game.UnitMap {
		jgame.UnitMap[fmt.Sprintf("%d,%d", loc.X, loc.Y)] = unit
	}
	return json.Marshal(jgame)
}

func (game *Game) UnmarshalJSON(data []byte) error {
	fmt.Println("Using the Unmarshal Function")
	jgame := &jsonGame{}
	err := json.Unmarshal(data, jgame)
	if err != nil {
		return err
	}
	game.Turn = jgame.Turn
	game.PlayerStatus = jgame.PlayerStatus
	game.UnitMap = map[Location]Unit{}
	for key, unit := range jgame.UnitMap {
		loc := Location{}
		fmt.Sscanf(key, "%d,%d", &loc.X, &loc.Y)
		game.UnitMap[loc] = unit
	}
	return nil
}
