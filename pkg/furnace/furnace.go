package furnace

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"

	api "github.com/prmsrswt/foundry/pkg/furnace/furnacepb"
)

type Furnace struct {
	mx sync.RWMutex
	// TODO(prmsrswt): Implement better queueing system. Something like https://github.com/sheerun/queue.
	buildQueue     []Package
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

// Enqueue adds a package from the build queue.
func (f *Furnace) enqueue(packages []Package) {
	f.mx.Lock()
	defer f.mx.Unlock()
	f.buildQueue = append(f.buildQueue, packages...)
}

// Dequeue removes a package from the build queue.
func (f *Furnace) dequeue() Package {
	f.mx.Lock()
	defer f.mx.Unlock()
	pkg := f.buildQueue[0]
	f.buildQueue = f.buildQueue[1:]
	return pkg
}

func (f *Furnace) Build(ctx context.Context, req *api.BuildRequest) (*api.BuildResponse, error) {
	packages := make([]Package, len(req.GetPackages()))
	for i, v := range req.GetPackages() {
		packages[i] = Package{Name: v.GetName(), Version: v.GetVersion()}
	}
	// TODO(prmsrswt): Check if any on these packages is already queued.
	f.enqueue(packages)
	level.Info(f.logger).Log("rpc", "build", "msg", "enqueued packages", "packages", fmt.Sprint(packages))
	return &api.BuildResponse{}, nil
}

// TODO(prmsrswt): Implement IsQueued.

type Package struct {
	Name    string
	Version string
}
