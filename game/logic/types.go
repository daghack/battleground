package types // import "github.com/daghack/battlegrounds/game/logic"

import (
	"encoding/json"
	"fmt"
	"strings"
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
	Id          string `json:"unitId"`
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

func (bs UnitMap) MarshalJSON() ([]byte, error) {
	contents := []string{}
	for k, v := range bs {
		unitBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		contents = append(contents, fmt.Sprintf(`"%d,%d" : %s`, k.X, k.Y, string(unitBytes)))
	}
	return []byte("{" + strings.Join(contents, ", ") + "}"), nil
}

func (bs *UnitMap) UnmarshalJSON(data []byte) error {
	gamestate := *bs
	unitMap := map[string]Unit{}
	err := json.Unmarshal(data, &unitMap)
	if err != nil {
		return err
	}
	for k, v := range unitMap {
		loc := Location{}
		fmt.Sscanf(k, "%d,%d", &loc.X, &loc.Y)
		gamestate[loc] = v
	}
	return nil
}
