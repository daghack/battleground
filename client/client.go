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

func update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf(`GameId: "%s"`, gameId))
	return nil
}

func main() {
	if len(os.Args) > 2 && os.Args[1] == "create" {
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
	err := ebiten.Run(update, screenWidth, screenHeight, 2, "Battleground Client");
	if err != nil {
		panic(err)
	}
}