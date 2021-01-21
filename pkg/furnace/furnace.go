package furnace

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/mikkeloscar/aur"
	"google.golang.org/grpc"

	api "github.com/prmsrswt/foundry/pkg/furnace/furnacepb"
)

type Furnace struct {
	buildQueue     chan aur.Pkg
	pkgs           pkgMap
	maxConcurrency int
	logger         log.Logger
	api.UnimplementedFurnaceServer
}

func NewFurnace(maxConcurrency, queueLength int, logger log.Logger) Furnace {
	return Furnace{
		pkgs: pkgMap{
			pkgs: make(map[string]status),
		},
		buildQueue:     make(chan aur.Pkg, queueLength),
		maxConcurrency: maxConcurrency,
		logger:         logger,
	}
}

func RegisterFurnaceServer(s *grpc.Server, fc *Furnace) {
	api.RegisterFurnaceServer(s, fc)
}

func (f *Furnace) Build(ctx context.Context, req *api.BuildRequest) (*api.BuildResponse, error) {
	pkgNames := make([]string, len(req.GetPackages()))
	for i, v := range req.GetPackages() {
		pkgNames[i] = v.GetName()
	}
	packages, err := aur.Info(pkgNames)
	if err != nil {
		// TODO(prmsrswt): Better error handling (gRPC codes etc.)
		return nil, err
	}
	for _, p := range packages {
		err := f.pkgs.setQueued(p)
		if err != nil {
			// The package is already queued/building, skip adding to the queue.
			continue
		}
		f.buildQueue <- p
	}
	level.Info(f.logger).Log("rpc", "build", "msg", "enqueued packages", "packages", fmt.Sprint(pkgNames))
	return &api.BuildResponse{
		Message: "packages added to build queue",
	}, nil
}

func (f *Furnace) IsQueued(ctx context.Context, req *api.IsQueuedRequest) (*api.IsQueuedResponse, error) {
	sts := f.pkgs.get(aur.Pkg{Name: req.GetPackage().GetName()})
	return &api.IsQueuedResponse{
		Status: sts == StatusBuilding || sts == StatusQueued,
	}, nil
}

type Package struct {
	Name string
}
