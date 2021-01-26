package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/metalmatze/signal/healthcheck"
	"github.com/metalmatze/signal/internalserver"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

func main() {
	conf := globalConfig{}
	rootCmd := &cobra.Command{
		Use:   "foundry",
		Short: "Foundry is a custom repository for pacman.",
		Long:  "A custom repository for pacman that builds packages from AUR on-demand.",
	}
	conf.registerFlags(rootCmd)

	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestamp, "caller", log.DefaultCaller)

	metrics := prometheus.NewRegistry()
	metrics.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)

	var g run.Group
	registerFurnace(rootCmd, &g, log.With(logger, "component", "furnace"), metrics)

	var isHelp bool
	{
		defaultHelpFunc := rootCmd.HelpFunc()
		rootCmd.SetHelpFunc(func(c *cobra.Command, s []string) {
			isHelp = true
			defaultHelpFunc(c, s)
		})
	}
	if err := rootCmd.Execute(); err != nil {
		level.Error(logger).Log("err", err)
		return
	}
	// Short circuit in case help command is called.
	if isHelp {
		return
	}

	{
		sig := make(chan os.Signal, 1)
		g.Add(func() error {
			signal.Notify(sig, os.Interrupt, syscall.SIGINT)
			<-sig
			return nil
		}, func(_ error) {
			level.Info(logger).Log("msg", "caught interrupt, exiting")
			signal.Stop(sig)
			close(sig)
		})
	}
	{
		healthchecks := healthcheck.NewMetricsHandler(healthcheck.NewHandler(), metrics)
		srv := internalserver.NewHandler(
			internalserver.WithName("Internal - Foundry API"),
			internalserver.WithHealthchecks(healthchecks),
			internalserver.WithPrometheusRegistry(metrics),
			internalserver.WithPProf(),
		)
		s := http.Server{
			Addr:    conf.listenInternal,
			Handler: srv,
		}
		g.Add(func() error {
			level.Info(logger).Log("msg", "starting internal HTTP server", "address", conf.listenInternal)
			return s.ListenAndServe()
		}, func(_ error) {
			// TODO(prmsrswt): Replace hardcoded timeout with a flag.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			_ = s.Shutdown(ctx)
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("msg", err.Error())
	}
}

type globalConfig struct {
	listenInternal string
}

func (fc *globalConfig) registerFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&fc.listenInternal, "internal-address", ":10200", "The address on which internal server listens.")
}
