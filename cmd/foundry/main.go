package main

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)

	level.Info(logger).Log("msg", "starting foundry")
}
