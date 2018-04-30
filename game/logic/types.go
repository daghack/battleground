package types // import "github.com/daghack/battlegrounds/game/logic"

import (
	"fmt"
	"strings"
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

func (loc Location) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d,%d", loc.X, loc.Y)), nil
}

func (loc *Location) UnmarshalJSON(data []byte) {
	fmt.Println(string(data))
	str := string(data)
	fmt.Sscanf(str, "%d,%d", &loc.X, &loc.Y)
}
