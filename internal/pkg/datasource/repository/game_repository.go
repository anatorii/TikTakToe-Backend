package repository

import (
	"log"
	"tiktaktoe/internal/pkg/datasource"
	"tiktaktoe/internal/pkg/datasource/mapper"
	"tiktaktoe/internal/pkg/domain"
	"time"

	"github.com/google/uuid"
)

type gameRepository struct {
	Postgres datasource.Postgres
}

func NewGameRepository() gameRepository {
	return gameRepository{
		Postgres: datasource.NewPostgres(),
	}
}

func (r gameRepository) CreateGame(game domain.Game) error {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
        INSERT INTO games (id, board, turn, status, player1_id, player2_id, single)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) 
        DO UPDATE SET 
            board      = EXCLUDED.board,
            turn       = EXCLUDED.turn,
            status     = EXCLUDED.status,
            player1_id = EXCLUDED.player1_id,
            player2_id = EXCLUDED.player2_id,
            single     = EXCLUDED.single
    `
	ds := mapper.GameToDatasource(&game)
	boardj, _ := ds.BoardJson()
	err := pool.Exec(query,
		ds.Id,
		boardj,
		ds.Turn,
		ds.Status,
		ds.Player1_id,
		ds.Player2_id,
		ds.Single,
	)
	if err != nil {
		log.Fatalf("Query error: %v", err)
		return err
	}
	return nil
}

func (r gameRepository) GetGame(id uuid.UUID) (domain.Game, error) {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
		SELECT games.id, games.board, games.turn, games.status, games.player1_id, games.player2_id, games.single
		FROM games
		WHERE games.id = $1
	`
	var game datasource.Game
	row, err := pool.QueryRow(query, id.String())
	err = row.Scan(&game.Id, &game.Board, &game.Turn, &game.Status, &game.Player1_id, &game.Player2_id, &game.Single)
	if err != nil {
		return domain.Game{}, err
	}
	return *mapper.DatasourceToGame(&game), nil
}

func (r gameRepository) UpdateGame(game domain.Game) (domain.Game, error) {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
        INSERT INTO games (id, board, turn, status, player1_id, player2_id, single)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) 
        DO UPDATE SET 
            board      = EXCLUDED.board,
            turn       = EXCLUDED.turn,
            status     = EXCLUDED.status,
            player1_id = EXCLUDED.player1_id,
            player2_id = EXCLUDED.player2_id,
            single     = EXCLUDED.single
    `
	ds := mapper.GameToDatasource(&game)
	boardj, _ := ds.BoardJson()
	err := pool.Exec(query,
		ds.Id,
		boardj,
		ds.Turn,
		ds.Status,
		ds.Player1_id,
		ds.Player2_id,
		ds.Single,
	)
	if err != nil {
		log.Fatalf("Query error: %v", err)
		return domain.Game{}, err
	}
	return game, nil
}

func (r gameRepository) DeleteGame(id uuid.UUID) error {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
        DELETE FROM games WHERE id = $1
    `
	err := pool.Exec(query, id)
	if err != nil {
		log.Fatalf("Query error: %v", err)
		return err
	}
	return nil
}

func (r gameRepository) ListGames(user domain.User) ([]domain.Game, error) {
	pool := r.Postgres.Pool
	if err := pool.Connect(); err != nil {
		return nil, err
	}
	defer pool.Close()

	query := `
		SELECT games.id, games.status, games.turn, games.board,
		games.single, games.player1_Id, games.player2_Id,
		games.created_at
		FROM games
		WHERE status not in ('player1_win', 'player2_win', 'draw') AND
		((single = 1 AND (
				player1_id = $1 or
				player1_id = uuid_nil())
				) or
			(single = 2 AND (
				player1_id = $1 or
			 	player2_id = $1 or
			 	player1_id = uuid_nil() or
			  	player2_id = uuid_nil())
		))
	`

	rows, err := pool.Query(query, user.Id)
	if err != nil {
		return nil, err
	}
	var game domain.Game
	var id uuid.UUID
	var status string
	var single int
	var p1 uuid.UUID
	var p2 uuid.UUID
	var cd time.Time
	var board domain.Board
	var turn int
	gms := make([]domain.Game, 0)
	for rows.Next() {
		err = rows.Scan(&id, &status, &turn, &board, &single, &p1, &p2, &cd)
		if err != nil {
			break
		}
		game.SetId(id)
		game.SetStatus(status)
		game.SetSingle(single)
		game.SetPlayer(1, p1)
		game.SetPlayer(2, p2)
		game.SetCreatedAt(cd.Unix())
		game.SetActiveTurn(domain.PlayerChar(turn))
		game.SetBoard(board)
		gms = append(gms, game)
	}
	return gms, err
}

func (r gameRepository) FinishedGames(userId uuid.UUID) ([]domain.Game, error) {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
		SELECT games.id, games.status, games.turn, games.board,
		games.single, games.player1_Id, games.player2_Id,
		games.created_at
		FROM games
		WHERE status in ('player1_win', 'player2_win', 'draw') AND
      		(player1_Id = $1 OR player2_Id = $1)
	`
	rows, err := pool.Query(query, userId)
	if err != nil {
		return nil, err
	}
	var game domain.Game
	var id uuid.UUID
	var status string
	var single int
	var p1 uuid.UUID
	var p2 uuid.UUID
	var cd time.Time
	var board domain.Board
	var turn int
	gms := make([]domain.Game, 0)
	for rows.Next() {
		err = rows.Scan(&id, &status, &turn, &board, &single, &p1, &p2, &cd)
		if err != nil {
			break
		}
		game.SetId(id)
		game.SetStatus(status)
		game.SetSingle(single)
		game.SetPlayer(1, p1)
		game.SetPlayer(2, p2)
		game.SetCreatedAt(cd.Unix())
		game.SetActiveTurn(domain.PlayerChar(turn))
		game.SetBoard(board)
		gms = append(gms, game)
	}
	return gms, err
}

func (r gameRepository) TopPlayers(numPlayers int) ([]domain.Player, error) {
	pool := r.Postgres.Pool
	pool.Connect()
	defer pool.Close()

	query := `
		select
			player,
			count(*) filter (where status_game = 'win')  as wins,
			count(*) filter (where status_game = 'loses') as loses,
			count(*) filter (where status_game = 'draw') as draws,
			round((count(*) filter (where status_game = 'win')) * 1.0 / count(*) * 100) as percent
		from (
			select
				player1_id as player,
				case 
					when status = 'player1_win' then 'win' 
					when status = 'player2_win' then 'loses'
					else 'draw'
				end as status_game
			from games
			where status in ('player1_win', 'player2_win', 'draw')
			and player1_id <> '00000000-0000-0000-0000-000000000000'::uuid

			union all

			select
				player2_id as player,
				case 
					when status = 'player2_win' then 'win' 
					when status = 'player1_win' then 'loses'
					else 'draw'
				end as status_game
			from games
			where status in ('player1_win', 'player2_win', 'draw')
			and player2_id <> '00000000-0000-0000-0000-000000000000'::uuid
		) t
		group by player
		order by 5 desc
		limit $1
	`
	rows, err := pool.Query(query, numPlayers)
	if err != nil {
		return nil, err
	}
	var player domain.Player
	var percent int
	ps := make([]domain.Player, 0)
	for rows.Next() {
		rows.Scan(&player.Id, &player.Wins, &player.Loses, &player.Draws, &percent)
		ps = append(ps, player)
	}
	return ps, nil
}
