package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/mikkeloscar/aur"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Builder interface {
	Build(aur.Pkg) error
}

// NoOpBuilder implements Builder interface. It's just a stub
// used for testing.
type NoOpBuilder struct {
	ShouldFail bool
}

func (nb *NoOpBuilder) Build(aur.Pkg) error {
	if nb.ShouldFail {
		return fmt.Errorf("build failed")
	}
	return nil
}

// MakepkgBuilder implements the Builder interface. It builds the
// packages using the `makepkg` script.
type MakepkgBuilder struct {
	wd string
}

func NewMakepkgBuilder(workDir string) *MakepkgBuilder {
	return &MakepkgBuilder{
		wd: workDir,
	}
}

func (mb *MakepkgBuilder) build(pkg aur.Pkg) error {
	pkgDir := path.Join(mb.wd, pkg.Name)
	repoURL := "https://aur.archlinux.org/" + pkg.Name + ".git"
	if err := gitClone(pkgDir, repoURL); err != nil {
		return err
	}
	cmd := exec.Command("makepkg", "-dfmc")
	cmd.Dir = pkgDir
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Build takes a package as input and starts building it locally.
// This requires `makepkg` to be present in the PATH to work.
// TODO(prmsrswt): Remove dependency on `makepkg` by doing
// everything natively.
func (mb *MakepkgBuilder) Build(pkg aur.Pkg) error {
	if err := mb.build(pkg); err != nil {
		return fmt.Errorf("building package: %w", err)
	}
	return nil
}

// gitClone will do a fresh clone, always. If the directory exists, it removes
// it completely before cloning.
func gitClone(dir, url string) error {
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if err = os.RemoveAll(dir); err != nil {
			return err
		}
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{URL: url, Depth: 1})
	return err
}

type instrumentedBuilder struct {
	builder Builder

	buildOps          prometheus.Counter
	buildOpsCompleted prometheus.Counter
	buildOpsFailure   prometheus.Counter

	buildDuration prometheus.Histogram
}

func BuilderWithMetrics(b Builder, reg prometheus.Registerer) *instrumentedBuilder {
	return &instrumentedBuilder{
		builder: b,
		buildOps: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "foundry_furnace_builder_builds_started_total",
			Help: "Total number of attempted builds by the builder.",
		}),
		buildOpsCompleted: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "foundry_furnace_builder_builds_completed_total",
			Help: "Total number of builds completed by the builder irrespective of success or failure.",
		}),
		buildOpsFailure: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "foundry_furnace_builder_build_failures_total",
			Help: "Total number of builds attempted by the builder that failed.",
		}),
		buildDuration: promauto.With(reg).NewHistogram(prometheus.HistogramOpts{
			Name:    "foundry_furnace_builder_build_duration_seconds",
			Help:    "Duration of builds completed successfully by the builder.",
			Buckets: []float64{10, 30, 60, 90, 120, 180, 300, 600, 900, 1800},
		}),
	}
}

func (ib *instrumentedBuilder) Build(pkg aur.Pkg) error {
	ib.buildOps.Inc()
	defer ib.buildOpsCompleted.Inc()
	start := time.Now()

	err := ib.builder.Build(pkg)
	if err != nil {
		ib.buildOpsFailure.Inc()
		return err
	}

	ib.buildDuration.Observe(time.Since(start).Seconds())
	return nil
}
