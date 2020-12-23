package main

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)

	rootCmd := &cobra.Command{
		Use:   "foundry",
		Short: "Foundry is a custom repository for pacman.",
		Long:  "A custom repository for pacman that builds packages from AUR on-demand.",
	}
	registerFurnace(rootCmd, log.With(logger, "component", "furnace"))

	if err := rootCmd.Execute(); err != nil {
		level.Error(logger).Log("msg", err.Error())
	}
}
