package main // import "github.com/daghack/battleground/server"

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	model "github.com/daghack/battleground/game/logic"
)

var games []*model.GameState

func init() {
	games = []*model.GameState{}
}

type createGameRequest struct {
	PlayerId model.Id `json:"playerId"`
}
type createGameResponse struct {
	GameId model.Id `json:"gameId"`
}

func createGame(w http.ResponseWriter, r *http.Request) {
	gameId := model.Id(fmt.Sprintf("%d", len(games)))
	jsonReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer r.Body.Close()

	args := createGameRequest{}
	err = json.Unmarshal(jsonReq, &args)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	resp := createGameResponse{ GameId : gameId }
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	games = append(games, model.NewGame(args.PlayerId))
	w.Write(jsonResp)
}

func joinGame(w http.ResponseWriter, r *http.Request) {
}

func readyPlayer(w http.ResponseWriter, r *http.Request) {
	type readyPlayerRequest struct {
		PlayerId model.Id `json:"playerId"`
		GameId model.Id `json:"gameId"`
		Field []model.UnitType `json:"field"`
	}
	type readyPlayerResponse struct {
		GameState *model.GameState `json:"gameState"`
	}
	jsonReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer r.Body.Close()

	args := readyPlayerRequest{}
	err = json.Unmarshal(jsonReq, &args)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	gameId := 0
	fmt.Sscanf(string(args.GameId), "%d", &gameId)
	games[gameId].ReadyPlayer(args.PlayerId, args.Field)
	resp := readyPlayerResponse{ GameState : games[gameId] }
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	w.Write(jsonResp)
}

func takeTurn(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/create_game", createGame)
	http.HandleFunc("/join_game", joinGame)
	http.HandleFunc("/ready_player", readyPlayer)
	http.HandleFunc("/take_turn", takeTurn)
	http.ListenAndServe(":6969", nil)
}
