package main

import (
	"net"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/prmsrswt/foundry/pkg/furnace"
)

func registerFurnace(cmd *cobra.Command, logger log.Logger) {
	config := &furnaceConfig{}
	furnaceCmd := &cobra.Command{
		Use:   "furnace",
		Short: "Run the Furnace component",
		Run: func(cmd *cobra.Command, args []string) {
			runFurnace(config, logger)
		},
	}
	cmd.AddCommand(furnaceCmd)
	config.registerFlags(furnaceCmd)
}

func runFurnace(config *furnaceConfig, logger log.Logger) {
	conn, err := net.Listen("tcp", config.grpc.bindAddress)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
	}
	s := grpc.NewServer()
	fc := furnace.NewFurnace(config.maxConcurrency, logger)
	furnace.RegisterFurnaceServer(s, &fc)
	level.Info(logger).Log("msg", "starting gRPC server", "addr", config.grpc.bindAddress)
	if err = s.Serve(conn); err != nil {
		level.Error(logger).Log("msg", err.Error())
	}
}

type furnaceConfig struct {
	grpc           grpcConfig
	maxConcurrency int
}

func (fc *furnaceConfig) registerFlags(cmd *cobra.Command) {
	fc.grpc.registerFlag(cmd)
	cmd.Flags().IntVar(&fc.maxConcurrency, "max-concurrency", 1, "Maximum number of packages to build at the same time")
}