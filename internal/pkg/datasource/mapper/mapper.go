package mapper

import (
	"tiktaktoe/internal/pkg/datasource"
	"tiktaktoe/internal/pkg/domain"
)

func GameToDatasource(m *domain.Game) *datasource.Game {
	return &datasource.Game{
		Id:         m.GetId(),
		Board:      m.GetBoard(),
		Turn:       int(m.GetActiveTurn()),
		Status:     m.GetStatus(),
		Player1_id: m.GetPlayer(1),
		Player2_id: m.GetPlayer(2),
		Single:     m.GetSingle(),
	}
}

func DatasourceToGame(ds *datasource.Game) *domain.Game {
	g := *domain.NewGame()
	g.SetId(ds.Id)
	g.SetBoard(ds.Board)
	g.SetActiveTurn(domain.PlayerChar(ds.Turn))
	g.SetStatus(ds.Status)
	g.SetPlayer(1, ds.Player1_id)
	g.SetPlayer(2, ds.Player2_id)
	g.SetSingle(ds.Single)
	return &g
}

func UserToDatasource(m *domain.User) *datasource.User {
	return &datasource.User{
		Id:       m.Id,
		Login:    m.Login,
		Password: m.Password,
	}
}

func DatasourceToUser(ds *datasource.User) *domain.User {
	return &domain.User{
		Id:       ds.Id,
		Login:    ds.Login,
		Password: ds.Password,
	}
}
