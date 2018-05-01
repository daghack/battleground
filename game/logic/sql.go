package types // import "github.com/daghack/battlegrounds/game/logic"

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
)

const fetchPlayer string = `select * from players where id=$1`
const createPlayer string = `insert into players (username, passkey) values (:username, :passkey) returning id`
const verifyPlayer string = `select id from players where username=$1 and passkey=$2`

const fetchGame string = `select * from active_games where id=$1`
const createGame string = `insert into active_games (board_size, piece_count, game_state) values (:board_size, :piece_count, :game_state) returning id`
const joinGame string = `insert into active_players (game_id, player_id) values ($1, $2)`
const updateGame string = `update active_games set game_state=$2 where id=$1`
const playersInGame string = `select player_id from active_players where game_id=$1`
const removeGame string = `delete from active_games where id=$1`
const leaveGame string = `delete from active_players where game_id=$1 and player_id=$2`

type DBHandler struct {
	db *sqlx.DB
}

func loadDatabase(postgresStr, loadfile string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", postgresStr)
	if err != nil {
		return nil, err
	}
	fileH, err := os.Open(loadfile)
	if err != nil {
		return nil, err
	}
	defer fileH.Close()
	schema, err := ioutil.ReadAll(fileH)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(schema))
	return db, err
}

func NewDBHandler(postgresStr, loadfile string) (*DBHandler, error) {
	db, err := loadDatabase(postgresStr, loadfile)
	if err != nil {
		return nil, err
	}
	return &DBHandler{db: db}, nil
}

func (dbh *DBHandler) FetchPlayer(playerId string) (*Player, error) {
	player := &Player{}
	err := dbh.db.Get(player, fetchPlayer, playerId)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (dbh *DBHandler) VerifyPlayer(username, passkey string) (string, error) {
	// Obviously, in a real world application, this behavior is blatantly incorrect,
	// as there is nothing that stops a user from creating an account with an identical
	// username and passkey. Even so, it is *good enough* as long as nobody ever plays it.
	// Consider this an area that needs to be fixed the second there are more than 2
	// players in the whole world.
	playerId string
	err := dbh.db.Select(&playerId, verifyPlayer, username, passkey)
	return playerId, err
}

func (dbh *DBHandler) CreatePlayer(player *Player) (string, error) {
	rows, err := dbh.db.NamedQuery(createPlayer, player)
	if err != nil {
		return "", err
	}
	id := ""
	for rows.Next() {
		err = rows.Scan(&id)
	}
	return id, err
}

func (dbh *DBHandler) FetchGame(gameId string) (*ActiveGame, error) {
	activeGame := &ActiveGame{}
	err := dbh.db.Get(activeGame, fetchGame, gameId)
	if err != nil {
		return nil, err
	}
	return activeGame, nil
}

func (dbh *DBHandler) CreateGame(boardSize, pieceCount int) (string, error) {
	game := &ActiveGame{
		BoardSize:  boardSize,
		PieceCount: pieceCount,
		GameState:  []byte(`{}`),
	}
	rows, err := dbh.db.NamedQuery(createGame, game)
	if err != nil {
		return "", err
	}
	id := ""
	for rows.Next() {
		err = rows.Scan(&id)
	}
	return id, err
}

func (dbh *DBHandler) JoinGame(gameId, playerId string) error {
	_, err := dbh.db.Exec(joinGame, gameId, playerId)
	return err
}

func (dbh *DBHandler) UpdateGame(gameId string, gamestate *Game) error {
	gamebytes, err := json.Marshal(gamestate)
	if err != nil {
		return err
	}
	_, err = dbh.db.Exec(updateGame, gameId, gamebytes)
	return err
}

func (dbh *DBHandler) ActivePlayersInGame(gameId string) ([]string, error) {
	playerIds := []string{}
	err := dbh.db.Select(&playerIds, playersInGame, gameId)
	return playerIds, err
}

func (dbh *DBHandler) RemoveGame(gameId string) error {
	_, err := dbh.db.Exec(removeGame, gameId)
	return err
}

func (dbh *DBHandler) LeaveGame(gameId, playerId string) error {
	_, err := dbh.db.Exec(leaveGame, gameId, playerId)
	return err
}
