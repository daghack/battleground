package main // import "github.com/daghack/battleground/server"

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	model "github.com/daghack/battleground/game/logic"
)

//var games []*model.GameState
var games map[model.Id] *model.GameState

func init() {
	//games = []*model.GameState{}
  games = map[model.Id] *model.GameState{}
}

type createGameRequest struct {
	PlayerId model.Id `json:"playerId"`
}

type createGameResponse struct {
	GameId model.Id `json:"gameId"`
}

type JoinGameInput struct {
  PlayerId model.Id `json:"playerId"`
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

	resp := createGameResponse{
    GameId: gameId,
  }
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

  games[gameId] = model.NewGame(args.PlayerId)
	w.Write(jsonResp)
}

func joinGame(resp http.ResponseWriter, req *http.Request) {
  jsonInput, err := ioutil.ReadAll(req.Body)
  if err != nil {
    http.Error(resp, err.Error(), 400)
    return
  }

  var input JoinGameInput
  err = json.Unmarshal(jsonInput, &input)
  if err != nil {
    http.Error(resp, err.Error(), 400)
    return
  }

  // First, lookup the game.
  var game *model.GameState

  game, ok := games[input.GameId]
  if !ok {
    http.Error(resp, "Game not found", 404)
    return
  }

  err = game.AddPlayer(input.PlayerId)
  if err != nil {
    http.Error(resp, err.Error(), 400)
    return
  }

  http.Error(resp, "", 201)
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

	games[args.GameId].ReadyPlayer(args.PlayerId, args.Field)
	resp := readyPlayerResponse{ GameState : games[args.GameId] }
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
