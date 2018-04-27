package types // import "github.com/daghack/battlegrounds/game/types"

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
}

type GameState map[Location]Unit
