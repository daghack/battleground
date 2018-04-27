package main // import "github.com/daghack/battleground/server"

import (
	"fmt"
	"encoding/json"
	model "github.com/daghack/battlegrounds/game/types"
)

func main() {
	gs := make(GameState)
	bytes, err := json.Marshal(gs)
	fmt.Println(string(bytes))
}
