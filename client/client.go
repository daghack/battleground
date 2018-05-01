package main // import "github.com/daghack/battleground/client"

import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	model "github.com/daghack/battleground/game/logic"
	"image/color"
)

const (
	screenWidth = 400
	screenHeight = 400
	boardWidth = 8
	boardHeight = 8
	tilesize = 50
)

var mouseUp = true
var target = 0
var targets = [2]model.Location{}
var gameId = ""
var playerId = ""
var gamestate *model.Game
var colorBlue color.RGBA = color.RGBA{ R : 0, G : 0, B : 255, A : 0 }

func update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}
	for key, _ := range gamestate.UnitMap {
		ebitenutil.DrawRect(screen, float64(key.X * tilesize), float64(key.Y * tilesize), tilesize, tilesize, colorBlue)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && mouseUp {
		mx, my := ebiten.CursorPosition()
		tx := mx / tilesize
		ty := my / tilesize
		if tx < 0 || tx >= boardWidth {
			return nil
		}
		if ty < 0 || ty >= boardHeight {
			return nil
		}
		targets[target] = model.Location{X : tx, Y : ty}
		target += 1
		target %= 2
		mouseUp = false
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseUp = true
	}
	ebitenutil.DrawRect(screen, float64(targets[0].X * tilesize), float64(targets[0].Y * tilesize), tilesize, tilesize, color.White)
	ebitenutil.DrawRect(screen, float64(targets[1].X * tilesize), float64(targets[1].Y * tilesize), tilesize, tilesize, color.White)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("GameId: '%s'\nGameState: '%v'", gameId, gamestate))
	return nil
}

func init() {
	gamestate = &model.Game{
		Turn : 1,
		UnitMap : map[model.Location]model.Unit{model.Location{X:1, Y:1}: model.Unit{}},
		PlayerStatus : map[string]int{"Hello" : 5},
	}
}

func main() {
	if len(os.Args) > 2 && os.Args[1] == "create" {
		playerId = os.Args[2]
		createRequest := struct { PlayerId string `json:"playerId"` }{ PlayerId : os.Args[2] }
		reqJson, err := json.Marshal(createRequest)
		if err != nil {
			panic(err)
		}
		buff := bytes.NewBuffer(reqJson)
		req, err := http.NewRequest("GET", "http://localhost:6969/create_game", buff)
		if err != nil {
			panic(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		respJson, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		
		createResp := struct {GameId string `json:"gameId"`}{}
		err = json.Unmarshal(respJson, &createResp)
		gameId = createResp.GameId
	} else if len(os.Args) > 3 && os.Args[1] == "join" {
		playerId = os.Args[2]
		createRequest := struct {
			PlayerId string `json:"playerId"`
			GameId string `json:"gameId"`
		}{ PlayerId : os.Args[2], GameId : os.Args[3] }
		reqJson, err := json.Marshal(createRequest)
		if err != nil {
			panic(err)
		}
		buff := bytes.NewBuffer(reqJson)
		req, err := http.NewRequest("GET", "http://localhost:6969/join_game", buff)
		if err != nil {
			panic(err)
		}
		client := &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			panic(err)
		}
		gameId = os.Args[3]
	}
	readyRequest := struct {
		PlayerId string `json:"playerId"`
		GameId string `json:"gameId"`
		Field []string `json:"field"`
	}{PlayerId : playerId, GameId : gameId, Field: []string{"footman", "footman", "footman", "footman", "footman", "footman", "footman", "footman"}}
	fmt.Println("Submitting Ready Player")
	reqJson, err := json.Marshal(readyRequest)
	if err != nil {
		panic(err)
	}
	buff := bytes.NewBuffer(reqJson)
	req, err := http.NewRequest("GET", "http://localhost:6969/ready_player", buff)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("Reading Body")
	respJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(string(respJson))
	unmarshalTarget := struct{ Game *model.Game `json:"gameState"` }{Game : gamestate}
	err = json.Unmarshal(respJson, &unmarshalTarget)
	if err != nil {
		panic(err)
	}
	err = ebiten.Run(update, screenWidth, screenHeight, 2, "Battleground Client");
	if err != nil {
		panic(err)
	}
}
