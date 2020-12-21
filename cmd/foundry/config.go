package main

import "github.com/spf13/cobra"

type grpcConfig struct {
	bindAddress string
}

func (gc *grpcConfig) registerFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&gc.bindAddress,
		"grpc-address",
		"0.0.0.0:10201",
		"Listen ip:port address for gRPC endpoints. Make sure this address is routable from other components.",
	)
}
