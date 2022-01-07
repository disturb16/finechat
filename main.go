package main

import (
	"context"
	"log"
	"sync"

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

		fx.Invoke(
			configureLifeCycle,
		),
	)

	app.Run()
}

func dbConnection(ctx context.Context, cfg *configuration.Configuration) (*sqlx.DB, error) {
	return database.CreateMysqlConnection(ctx, cfg.DB)
}

func configureLifeCycle(ctx context.Context, lc fx.Lifecycle, db *sqlx.DB, s auth.Service) {
	lc.Append(fx.Hook{
		OnStart: func(ctxStart context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			wg := &sync.WaitGroup{}
			wg.Add(1)

			go func() {
				defer wg.Done()
				log.Println("Closing database connection...")
				err := db.Close()
				if err != nil {
					log.Println("Error closing database connection:", err)
				}
			}()

			wg.Wait()
			return nil
		},
	})
}
