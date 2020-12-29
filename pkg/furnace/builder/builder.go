package builder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/mikkeloscar/aur"
)

type status string

const (
	StatusBuilding status = "building"
	StatusError    status = "error"
	StatusDefault  status = ""
)

var (
	ErrAlreadyBuilding error = errors.New("package is already building")
)

type Builder interface {
	Build(aur.Pkg) error
}

// MakepkgBuilder implements the Builder interface. It builds the
// packages using the `makepkg` script.
type MakepkgBuilder struct {
	pkgStatus pkgMap
	wd        string
}

func NewMakepkgBuilder(workDir string) *MakepkgBuilder {
	return &MakepkgBuilder{
		wd: workDir,
		pkgStatus: pkgMap{
			pkgs: make(map[string]status),
		},
	}
}

func (mb *MakepkgBuilder) build(pkg aur.Pkg) error {
	pkgDir := path.Join(mb.wd, pkg.Name)
	repoURL := "https://aur.archlinux.org/" + pkg.Name + ".git"
	if err := gitClone(pkgDir, repoURL); err != nil {
		return err
	}
	cmd := exec.Command("makepkg", "-sf")
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
	err := mb.pkgStatus.setBuilding(pkg)
	if err != nil {
		return err
	}
	if err := mb.build(pkg); err != nil {
		mb.pkgStatus.set(pkg, StatusError)
		return fmt.Errorf("building package: %w", err)
	}
	mb.pkgStatus.set(pkg, StatusDefault)
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
