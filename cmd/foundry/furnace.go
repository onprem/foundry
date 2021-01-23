package main

import (
	"net"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/prmsrswt/foundry/pkg/furnace"
	"github.com/prmsrswt/foundry/pkg/furnace/builder"
)

func registerFurnace(cmd *cobra.Command, g *run.Group, logger log.Logger, metrics *prometheus.Registry) {
	config := &furnaceConfig{}
	furnaceCmd := &cobra.Command{
		Use:   "furnace",
		Short: "Run the Furnace component",
		Run: func(cmd *cobra.Command, args []string) {
			setupFurnace(config, g, logger, metrics)
		},
	}
	cmd.AddCommand(furnaceCmd)
	config.registerFlags(furnaceCmd)
}

func setupFurnace(config *furnaceConfig, g *run.Group, logger log.Logger, _ *prometheus.Registry) {
	fc := furnace.NewFurnace(config.maxConcurrency, config.queueLimit, logger)

	{
		conn, err := net.Listen("tcp", config.grpc.bindAddress)
		if err != nil {
			// TODO(prmsrswt): this is a non-recoverable error, handle it like one.
			level.Error(logger).Log("msg", err.Error())
		}
		s := grpc.NewServer()
		furnace.RegisterFurnaceServer(s, &fc)

		g.Add(func() error {
			level.Info(logger).Log("msg", "starting gRPC server", "addr", config.grpc.bindAddress)
			return s.Serve(conn)
		}, func(_ error) {
			s.GracefulStop()
		})
	}

	{
		makepkgBuilder := builder.NewMakepkgBuilder("/tmp/foundry/furnace")
		g.Add(func() error {
			fc.Start(makepkgBuilder)
			return nil
		}, func(_ error) {
			// TODO(prmsrswt): Use context canceling in fc.Start.
		})
	}
}

type furnaceConfig struct {
	grpc           grpcConfig
	maxConcurrency int
	queueLimit     int
}

func (fc *furnaceConfig) registerFlags(cmd *cobra.Command) {
	fc.grpc.registerFlag(cmd)
	cmd.Flags().IntVar(&fc.maxConcurrency, "max-concurrency", 1, "Maximum number of packages to build at the same time")
	cmd.Flags().IntVar(&fc.queueLimit, "queue-limit", 100, "Maximum number of packages to have in build queue at one time. If this limit is reached then the request will wait for the queue to accommodate the packages.")
}
