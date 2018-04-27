package types // import "github.com/daghack/battlegrounds/game/types"

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	NORTH int = iota
	SOUTH
	EAST
	WEST
)

type Location struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Orientation int

type Id string

type UnitType string

type Unit struct {
	Id Id `json:"unitId"`
	Orientation Orientation `json:"orientation"`
	UnitType UnitType `json:"unitType"`
	PlayerId Id `json:"playerId"`
}

type BoardState map[Location]Unit

type PlayerState struct {
	Id Id `json:"playerId"`
	Ready bool `json:"ready"`
	Field []UnitType `json:"field"`
}

type GameState struct {
	Id Id `json:"id"`
	Players map[Id]PlayerState `json:"players"`
	CurrentPlayer Id `json:"currentId"`
	CurrentTurn int `json:"currentTurn"`
}

func (bs BoardState) MarshalJSON() ([]byte, error) {
	contents := []string{}
	for k, v := range bs {
		unitBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		contents = append(contents, fmt.Sprintf(`"%d,%d" : %s`, k.X, k.Y, string(unitBytes)))
	}
	return []byte("{"+strings.Join(contents, ", ")+"}"), nil
}

func (bs *BoardState) UnmarshalJSON(data []byte) error {
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
