package caddyphpfpm

import (
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
)

// L is the global logger to be used by the app
// So when we need to log we import it and use log.L.<method>(...)
var L *zap.Logger

func init() {
	var err error
	L, err = zapdriver.NewProduction()
	if err != nil {
		panic(err)
	}
}
