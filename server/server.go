package main // import "github.com/daghack/battleground/server"

import (
	"fmt"
	"encoding/json"
	model "github.com/daghack/battleground/game/logic"
)

func main() {
	gs := model.GameState{Players : map[model.Id]model.PlayerState{"me" : model.PlayerState{}}}
	bytes, err := json.Marshal(gs)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	gs2 := model.GameState{}
	err = json.Unmarshal(bytes, &gs2)
	if err != nil {
		panic(err)
	}
	fmt.Println(gs2)
}
