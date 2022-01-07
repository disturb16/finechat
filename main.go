package main

import (
	"context"

	"github.com/disturb16/finechat/configuration"
	"github.com/disturb16/finechat/database"
	"github.com/disturb16/finechat/internal/auth"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			context.Background,
			configuration.GetDefault,
			dbConnection,
			auth.NewRepository,
			auth.NewService,
		),
	)

	app.Run()
}

func dbConnection(ctx context.Context, cfg *configuration.Configuration) (*sqlx.DB, error) {
	return database.CreateMysqlConnection(ctx, cfg.DB)
}
