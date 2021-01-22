package furnace

import (
	"errors"
	"sync"

	"github.com/mikkeloscar/aur"
)

type status string

const (
	StatusQueued   status = "queued"
	StatusBuilding status = "building"
	StatusBuilt    status = "built"
	StatusError    status = "error"
	StatusDefault  status = ""
)

var (
	ErrAlreadyBuilding error = errors.New("package is already building")
)

type pkgMap struct {
	sync.RWMutex
	pkgs map[string]status
}

func (p *pkgMap) get(pkg aur.Pkg) status {
	p.RLock()
	defer p.RUnlock()
	return p.pkgs[genKey(pkg)]
}

func (p *pkgMap) set(pkg aur.Pkg, s status) {
	p.Lock()
	defer p.Unlock()
	p.pkgs[genKey(pkg)] = s
}

// setQueued is a separate function to make this operation atomic/goroutine safe.
func (p *pkgMap) setQueued(pkg aur.Pkg) error {
	p.Lock()
	defer p.Unlock()
	if sts := p.pkgs[genKey(pkg)]; sts == StatusQueued || sts == StatusBuilding {
		return ErrAlreadyBuilding
	}
	p.pkgs[genKey(pkg)] = StatusQueued
	return nil
}

func genKey(pkg aur.Pkg) string {
	return pkg.Name
}
