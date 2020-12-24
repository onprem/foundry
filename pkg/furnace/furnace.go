package furnace

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"

	api "github.com/prmsrswt/foundry/pkg/furnace/furnacepb"
)

type Furnace struct {
	buildQueue     queue
	maxConcurrency int
	logger         log.Logger
	api.UnimplementedFurnaceServer
}

func NewFurnace(maxConcurrency int, logger log.Logger) Furnace {
	return Furnace{
		maxConcurrency: maxConcurrency,
		logger:         logger,
	}
}

func RegisterFurnaceServer(s *grpc.Server, fc *Furnace) {
	api.RegisterFurnaceServer(s, fc)
}

func (f *Furnace) Build(ctx context.Context, req *api.BuildRequest) (*api.BuildResponse, error) {
	packages := make([]Package, len(req.GetPackages()))
	for i, v := range req.GetPackages() {
		packages[i] = Package{Name: v.GetName()}
	}
	// TODO(prmsrswt): Check if any on these packages is already queued.
	f.buildQueue.enqueue(packages)
	level.Info(f.logger).Log("rpc", "build", "msg", "enqueued packages", "packages", fmt.Sprint(packages))
	return &api.BuildResponse{
		Message: "packages added to build queue",
	}, nil
}

func (f *Furnace) IsQueued(ctx context.Context, req *api.IsQueuedRequest) (*api.IsQueuedResponse, error) {
	return &api.IsQueuedResponse{
		Status: f.buildQueue.isQueued(Package{Name: req.GetPackage().GetName()}),
	}, nil
}

type Package struct {
	Name string
}
