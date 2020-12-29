package builder

import (
	"sync"

	"github.com/mikkeloscar/aur"
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

// setBuilding is a separate function to make this operation atomic/goroutine safe.
func (p *pkgMap) setBuilding(pkg aur.Pkg) error {
	p.Lock()
	defer p.Unlock()
	if p.pkgs[genKey(pkg)] == StatusBuilding {
		return ErrAlreadyBuilding
	}
	p.pkgs[genKey(pkg)] = StatusBuilding
	return nil
}

func genKey(pkg aur.Pkg) string {
	return pkg.Name
}
