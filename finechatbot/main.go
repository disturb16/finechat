package main

import (
	"context"
	"log"
	"net/http"

	"github.com/disturb16/finechatbot/bot"
	"github.com/disturb16/finechatbot/broker"
	"github.com/disturb16/finechatbot/configuration"
	"github.com/disturb16/finechatbot/healthcheckserver"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			configuration.GetDefault,
			func(config *configuration.Configuration) (broker.MessageBroker, error) {
				return broker.New(config)
			},
			bot.New,
			healthcheckserver.New,
		),
		fx.Invoke(
			setLyfeCycle,
		),
	)

	app.Run()
}

func setLyfeCycle(
	lc fx.Lifecycle,
	messageBroker broker.MessageBroker,
	b *bot.Bot,
	server *http.Server,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil {
					log.Fatal("FAILED TO START SERVER ", err)
				}
			}()

			go func() {
				err := b.Listen()
				if err != nil {
					log.Fatal("FAILED TO START BOT ", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			go func() {
				if err := server.Shutdown(ctx); err != nil {
					log.Fatal("FAILED TO STOP SERVER ", err)
				}
			}()

			go func() {
				err := messageBroker.Close()
				if err != nil {
					log.Fatal("FAILED TO CLOSE MESSAGE BROKER ", err)
				}
			}()

			return nil
		},
	})
}
