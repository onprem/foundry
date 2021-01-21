package furnace

import (
	"sync"

	"github.com/go-kit/kit/log/level"

	"github.com/prmsrswt/foundry/pkg/furnace/builder"
)

// TODO(prmsrswt): Gracefully clean this up on shutdown.
func (f *Furnace) Start(bldr builder.Builder) {
	level.Info(f.logger).Log("msg", "started package builder")
	wg := sync.WaitGroup{}
	for i := 0; i < f.maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			for pkg := range f.buildQueue {
				{
					f.pkgs.set(pkg, StatusBuilding)
					level.Debug(f.logger).Log("msg", "starting package build", "package", pkg.Name, "version", pkg.Version)
				}
				err := bldr.Build(pkg)
				if err != nil {
					f.pkgs.set(pkg, StatusError)
					level.Error(f.logger).Log("err", err, "package", pkg.Name)
					continue
				}
				{
					f.pkgs.set(pkg, StatusBuilt)
					level.Info(f.logger).Log("msg", "successfully built package", "package", pkg.Name)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
