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

func (gs *GameState) MarshalJSON() ([]byte, error) {
	contents := []string{}
	for k, v := range gs {
		unitBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		contents = append(contents, fmt.Sprintf(`"%d,%d" : %s`, k.X, k.Y, string(unitBytes)))
	}
	return []byte("{"+strings.Join(contents, ", ")+"}")

}

func (gs *GameState) UnmarshalJSON() ([]byte, error) {
	return nil, nil
}
