package logger

import (
	"context"
	"fmt"
	"log"

	"github.com/disturb16/finechat/configuration"
)

const (
	RequestIDKey string = "x-request-id"
	RealIPKey    string = "x-real-ip"
)

var debug bool

func Setup(conif *configuration.Configuration) {
	debug = conif.App.Debug
}

func contextArgs(ctx context.Context) []interface{} {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		requestID = "-"
	}

	realIP, ok := ctx.Value(RealIPKey).(string)
	if !ok {
		realIP = "-"
	}

	msgArgs := make([]interface{}, 0, 2)

	reqIDArg := fmt.Sprintf(", %s: %s", RequestIDKey, requestID)
	realIPArg := fmt.Sprintf(", %s: %s", RealIPKey, realIP)

	return append(msgArgs, reqIDArg, realIPArg)
}

func Println(ctx context.Context, args ...interface{}) {
	msgArgs := contextArgs(ctx)
	args = append(args, msgArgs...)
	log.Println(args...)
}
