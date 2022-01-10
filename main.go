package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/disturb16/finechat/broker"
	"github.com/disturb16/finechat/configuration"
	"github.com/disturb16/finechat/database"
	"github.com/disturb16/finechat/internal/api"
	"github.com/disturb16/finechat/internal/auth"
	"github.com/disturb16/finechat/internal/finechatbot"

	"github.com/disturb16/finechat/internal/chatroom"
	"github.com/disturb16/finechat/internal/client"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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
			chatroom.NewRepository,
			chatroom.NewService,
			echo.New,
			broker.New,
			api.NewHandler,
		),

		fx.Invoke(
			api.RegisterRoutes,
			configureLifeCycle,
			client.SetResources,
			func(b *broker.Broker) {
				go finechatbot.Listen(b)
			},
		),
	)

	app.Run()
}

func dbConnection(ctx context.Context, cfg *configuration.Configuration) (*sqlx.DB, error) {
	return database.CreateMysqlConnection(ctx, cfg.DB)
}

func configureLifeCycle(
	lc fx.Lifecycle,
	db *sqlx.DB,
	config *configuration.Configuration,
	e *echo.Echo,
	b *broker.Broker,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting server...")
			addr := fmt.Sprintf(":%d", config.App.Port)
			go func() {
				e.Logger.Fatal(e.Start(addr))
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			wg := &sync.WaitGroup{}
			wg.Add(3)

			go func() {
				defer wg.Done()
				log.Println("Closing database connection...")
				err := db.Close()
				if err != nil {
					log.Println("Error closing database connection:", err)
				}
			}()

			go func() {
				defer wg.Done()
				log.Println("Stopping server...")
				err := e.Shutdown(ctx)
				if err != nil {
					log.Println("Error shutting down server:", err)
				}
			}()

			go func() {
				defer wg.Done()
				log.Println("Stopping message broker...")
				err := b.Close()
				if err != nil {
					log.Println("Error shutting down message broker:", err)
				}
			}()

			wg.Wait()
			return nil
		},
	})
}
