package main // import "github.com/daghack/battleground/server"

import (
	"fmt"
	"encoding/json"
	model "github.com/daghack/battleground/game/types"
)

func main() {
	gs := make(model.GameState)
	gs[model.Location{X : 0, Y : 0}] = model.Unit{}
	bytes, err := json.Marshal(gs)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
