package main // import "github.com/daghack/battleground/server"

import (
	"encoding/json"
	model "github.com/daghack/battleground/game/logic"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
)

//var games []*model.GameState
var gameManager *model.GameManager

func init() {
	var err error
	gameManager, err = model.NewGameManager("user=postgres dbname=battleground sslmode=disable", "sql/state0.sql")
	if err != nil {
		panic(err)
	}
	playerId, err := gameManager.CreatePlayer("daghack", "supersecret")
	if err != nil {
		panic(err)
	}
	fmt.Println("Added User:", playerId)
}

type createGameRequest struct {
	PlayerId string `json:"playerId"`
}

type createGameResponse struct {
	GameId string `json:"gameId"`
}

type JoinGameInput struct {
	PlayerId string `json:"playerId"`
	GameId   string `json:"gameId"`
}

type MoveUnitInput struct {
	GameId        string         `json:"gameId"`
	PlayerId      string         `json:"playerId"`
	MoveFrom      model.Location `json:"moveFrom"`
	MoveTo        model.Location `json:"moveTo"`
	OrientTowards int            `json:"orientTowards"`
}

func createGame(resp http.ResponseWriter, r *http.Request) {
	jsonReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	args := createGameRequest{}
	err = json.Unmarshal(jsonReq, &args)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}

	log.Printf("Creating game %v", args)
	gameId, err := gameManager.CreateGame(args.PlayerId)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		log.Printf("Error creating game %s", err)
		return
	}

	cgResp := createGameResponse{
		GameId: gameId,
	}
	jsonResp, err := json.Marshal(cgResp)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		log.Printf("Error creating game %s", err)
		return
	}

	resp.Write(jsonResp)
	log.Printf("Created game %v", gameId)
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
	log.Printf("Joining game for %v", input)

	err = gameManager.JoinGame(input.GameId, input.PlayerId)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}
	http.Error(resp, "", 201)
	log.Printf("%v joined", input)
}

func readyPlayer(resp http.ResponseWriter, r *http.Request) {
	type readyPlayerRequest struct {
		PlayerId string   `json:"playerId"`
		GameId   string   `json:"gameId"`
		Field    []string `json:"field"`
	}
	type readyPlayerResponse struct {
		GameState *model.Game `json:"gameState"`
	}
	jsonReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	args := readyPlayerRequest{}
	err = json.Unmarshal(jsonReq, &args)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}

	err = gameManager.ReadyPlayer(args.GameId, args.PlayerId, args.Field)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}
	game, err := gameManager.FetchGame(args.GameId)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}
	rpResp := readyPlayerResponse{GameState: game}
	fmt.Println("Player Readied. Gamestate: ", *rpResp.GameState)
	jsonResp, err := json.Marshal(rpResp)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}
	resp.Write(jsonResp)
}

func moveUnit(resp http.ResponseWriter, req *http.Request) {
	jsonInput, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}

	var input MoveUnitInput
	err = json.Unmarshal(jsonInput, &input)
	if err != nil {
		http.Error(resp, err.Error(), 400)
		return
	}

	err = gameManager.MoveUnit(input.GameId, input.PlayerId,
		input.MoveFrom,
		input.MoveTo,
		input.OrientTowards,
	)
	if err != nil {
		http.Error(resp, err.Error(), 500)
		return
	}

	http.Error(resp, "", 203)
}

func main() {
	http.HandleFunc("/create_game", createGame)
	http.HandleFunc("/join_game", joinGame)
	http.HandleFunc("/ready_player", readyPlayer)
	http.HandleFunc("/move_unit", moveUnit)
	http.ListenAndServe(":6969", nil)
}
